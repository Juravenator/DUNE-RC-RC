package daq

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"text/template"
	"time"

	"cli.rc.ccm.dunescience.org/internal"
	"github.com/Masterminds/sprig"
)

// SendCommand sends a command and waits for return path response
func SendCommand(w internal.Writers, c *internal.RCConfig, name string, command string, timeout time.Duration) error {
	// fetch DAQ Application resource
	fmt.Fprintf(w.Out, "fetching daq application... ")
	app, err := internal.GetResource(c, internal.DAQAppKind, name)
	if err != nil {
		return err
	}
	fmt.Fprintln(w.Out, "OK")

	// find DAQ application location
	fmt.Fprintf(w.Out, "querying location of daq application... ")
	addr, err := internal.GetFirstServiceAddress(c, app.Spec["daq-service"].(string))
	if err != nil {
		return err
	}
	fmt.Fprintln(w.Out, "OK")

	// generate config
	fmt.Fprintf(w.Out, "generating config... ")
	payload, err := generateConfig(c, app.Spec, command)
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
		daqResponse, err := openReturnPath(listener, timeout)
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
	fmt.Fprintf(w.Out, "sending %s command to %s... ", command, name)
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
	var parsed []struct {
		Data interface{} `json:"data"`
		ID   string      `json:"id"`
	}
	err = json.Unmarshal(rendered.Bytes(), &parsed)
	if err != nil {
		return nil, err
	}
	for _, p := range parsed {
		if p.ID == command {
			return json.Marshal(p)
		}
	}
	return nil, fmt.Errorf("command '%s' not found in config", command)
}

type templateData struct {
	RunNumber string
}
