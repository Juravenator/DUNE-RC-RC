package daq

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"cli.rc.ccm.dunescience.org/internal"
	"github.com/Jeffail/gabs"
	"github.com/Masterminds/sprig"
)

// SendCommandOpts is an arg for SendCommand()
type SendCommandOpts struct {
	Rc       *internal.RCConfig
	DAQApp   string
	Command  string
	Timeout  time.Duration
	NewRunNr *uint64
}

// SendCommand sends a command and waits for return path response
func SendCommand(w internal.Writers, args SendCommandOpts) error {
	// fetch DAQ Application resource
	fmt.Fprintf(w.Out, "fetching daq application... ")
	app, err := internal.GetResource(args.Rc, internal.DAQAppKind, args.DAQApp)
	if err != nil {
		return err
	}
	fmt.Fprintln(w.Out, "OK")

	if isEnabled, exists := app.Spec["enabled"].(bool); exists && isEnabled {
		log.Error().Str("name", args.DAQApp).Msg("refusing to send command to DAQ application in autonomous mode")
		return fmt.Errorf("refusing to send command to DAQ application in autonomous mode")
	}

	// find DAQ application location
	fmt.Fprintf(w.Out, "querying location of daq application... ")
	addr, err := internal.GetFirstServiceAddress(args.Rc, app.Spec["daq-service"].(string))
	if err != nil {
		return err
	}
	fmt.Fprintln(w.Out, "OK")

	if args.NewRunNr != nil {
		app.Spec["run-number"] = fmt.Sprintf("%d", *args.NewRunNr)
		log.Warn().Uint64("new-run-number", *args.NewRunNr).Msg("setting new run number. This change is persistent, any commands following this one will use this new run number")
		err := internal.Apply(w, args.Rc, *app)
		if err != nil {
			log.Error().Err(err).Msg("cannot apply the new run number")
			return err
		}
	}

	// generate config
	fmt.Fprintf(w.Out, "generating config... ")
	payload, err := generateConfig(args.Rc, app.Spec, args.Command)
	if err != nil {
		return err
	}
	fmt.Fprintln(w.Out, "OK")
	log.Trace().Bytes("config", payload).Msg("rendered config")

	// setup return path http server
	fmt.Fprintf(w.Out, "setting up return path... ")
	listener, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		return err
	}
	_, returnPort, _ := net.SplitHostPort(listener.Addr().String())
	log.Debug().Str("port", returnPort).Msg("opened return path")
	daqResponseChan := make(chan *string)
	go func() {
		daqResponse, err := openReturnPath(listener, args.Timeout)
		if err != nil {
			log.Error().Err(err).Msg("return path failed")
		}
		if daqResponse == nil {
			log.Debug().Msg("openReturnPath() returned nil")
		} else {
			log.Debug().Str("daqResponse", *daqResponse).Msg("openReturnPath() returned")
		}
		daqResponseChan <- daqResponse
	}()
	fmt.Fprintln(w.Out, "OK")

	// send command
	fmt.Fprintf(w.Out, "sending %s command to %s... ", args.Command, args.DAQApp)
	url := fmt.Sprintf("http://%s:%d/command", addr.Address, addr.Port)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Answer-Port", returnPort)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("sending command failed with code %s", resp.Status)
	}

	fmt.Fprintln(w.Out, "OK")

	// wait for response
	fmt.Fprintf(w.Out, "waiting for return path... ")
	daqResponse := <-daqResponseChan
	if daqResponse != nil && *daqResponse == "OK" {
		fmt.Fprintln(w.Out, "OK")
		return nil
	}
	if daqResponse == nil {
		fmt.Fprintln(w.Out, "TIMEOUT")
		r := "command timed out"
		daqResponse = &r
	} else {
		fmt.Fprintln(w.Out, "FAIL")
	}
	log.Debug().Str("daqresponse", *daqResponse).Msg("daqresponse")
	return fmt.Errorf(*daqResponse)
}

func generateConfig(c *internal.RCConfig, spec internal.GenericSpec, command string) ([]byte, error) {
	configKey := spec["configkey"].(string)[len("/daq-configs/"):]
	resource, err := internal.GetResource(c, internal.DAQConfigKind, configKey)
	if err != nil {
		return nil, err
	}

	templated := resource.Spec["template"].(string)
	input := ""
	scanner := bufio.NewScanner(strings.NewReader(templated))
	for scanner.Scan() {
		line := scanner.Text()
		escapedLine := ""

		for line != "" {
			tmplStartI := strings.Index(line, "{{")
			if tmplStartI == -1 {
				escapedLine += line
				break
			}
			tmplEndI := strings.Index(line[tmplStartI:], "}}")
			if tmplEndI == -1 {
				escapedLine += line
				break
			}
			tmplEndI += tmplStartI
			escapedLine += line[:tmplStartI]
			escapedLine += strings.ReplaceAll(line[tmplStartI:tmplEndI+2], "\\\"", "\"")
			line = line[tmplEndI+2:]
		}

		input += escapedLine + "\n"
	}
	funcs := template.FuncMap{
		"firstServiceAddr": func(serviceName string) internal.ServiceAddress {
			addr, err := internal.GetFirstServiceAddress(c, serviceName)
			if err != nil {
				panic(err)
			}
			return *addr
		},
	}
	inputTemplate, err := template.New("").Funcs(sprig.TxtFuncMap()).Funcs(funcs).Parse(input)
	if err != nil {
		panic(err)
	}

	var rendered bytes.Buffer
	err = inputTemplate.Execute(&rendered, templateData{
		RunNumber: spec["run-number"].(string),
	})
	if err != nil {
		return nil, err
	}
	parsed, err := gabs.ParseJSON(rendered.Bytes())
	commandDatas, err := parsed.Children()
	if err != nil {
		return nil, fmt.Errorf("cannot parse DAQ config: %w", err)
	}
	runNr, err := strconv.ParseUint(spec["run-number"].(string), 10, 64)
	if err != nil {
		return nil, err
	}
	for _, commandData := range commandDatas {
		id := commandData.Search("id").Data().(string)
		if command == id {
			return injectRunNumber(commandData, runNr)
		}
	}

	return nil, fmt.Errorf("command '%s' not found in config", command)
}

type templateData struct {
	RunNumber string
}

func injectRunNumber(daqCommandConfig *gabs.Container, runNumber uint64) ([]byte, error) {
	commandModules, err := daqCommandConfig.Search("data", "modules").Children()
	if err != nil {
		return nil, fmt.Errorf("config item contains no modules: %w", err)
	}
	for _, module := range commandModules {
		if module.Exists("data", "run") {
			module.Set(runNumber, "data", "run")
		}
	}
	return daqCommandConfig.Bytes(), nil
}
