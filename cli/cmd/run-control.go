package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"cli.rc.ccm.dunescience.org/daq"
	"cli.rc.ccm.dunescience.org/internal"

	"github.com/rs/zerolog"
	logger "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

var log = logger.With().Str("pkg", "main").Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})

var rcConfig = internal.RCConfig{}
var writers = internal.Writers{}

var flags = []cli.Flag{
	&cli.StringFlag{Name: "config-file", Aliases: []string{"c"}, Usage: "location of a config file to use", EnvVars: []string{"RC_CONFIG"}},
	altsrc.NewStringFlag(&cli.StringFlag{Name: "log-level", Aliases: []string{"l"}, Value: "info"}),
	altsrc.NewStringFlag(&cli.StringFlag{Name: "rc-host", Aliases: []string{"host"}, Value: "localhost", Usage: "communicate with RC running on `hostname`", EnvVars: []string{"RC_HOST"}}),
	altsrc.NewUintFlag(&cli.UintFlag{Name: "consul-port", Value: 8500, Usage: "communicate with RC consul subsystem running on `port`", EnvVars: []string{"RC_CONSUL_PORT"}}),
	altsrc.NewUintFlag(&cli.UintFlag{Name: "nomad-port", Value: 4646, Usage: "communicate with RC nomad subsystem running on `port`", EnvVars: []string{"RC_NOMAD_PORT"}}),
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	app := &cli.App{
		Version:  "v0.0.0-march-rc",
		Compiled: time.Now(),
		Authors: []*cli.Author{{
			Name:  "Glenn Dirkx",
			Email: "glenn.dirkx@cern.ch",
		}},
		Copyright: "(c) 1989 DUNE",
		Usage:     "CLI for interfacing with a running RC cluster",
		Commands: []*cli.Command{
			{
				Name:        "get",
				Usage:       "display one or many resources",
				Description: "Prints a table of the most important information about the specified resource(s).\nUse 'run-control api-resources' for a complete list of supported resources.",
				// Flags: []cli.Flag{
				// 	&cli.StringFlag{Name: "output", Aliases: []string{"o"}},
				// },
				ArgsUsage: "kind [name]",
				Action: func(c *cli.Context) error {
					kind := c.Args().Get(0)
					id := c.Args().Get(1)
					if kind == "" {
						return fmt.Errorf("no kind given, try running 'run-control get all' or 'run-control api-resources'")
					}
					if kind == "all" {
						if id != "" {
							return fmt.Errorf("cannot supply ids when kind is 'all'")
						}
						// fetch and print all resources
						kinds, err := internal.GetAllKinds(&rcConfig)
						if err != nil {
							return err
						}
						fmt.Fprintln(c.App.Writer, "KIND\t\tNAME")
						for _, kind := range kinds {
							keys, err := internal.GetAllKeys(&rcConfig, kind)
							if err != nil {
								return err
							}
							for _, key := range keys {
								fmt.Fprintln(c.App.Writer, kind+"\t"+key)
							}
						}
						return nil
					}
					if id == "" {
						// fetch and print all resources of kind
						keys, err := internal.GetAllKeys(&rcConfig, kind)
						if err != nil {
							return err
						}
						fmt.Fprintln(c.App.Writer, "NAME")
						for _, key := range keys {
							fmt.Fprintln(c.App.Writer, key)
						}
						return nil
					}

					// fetch and print specific resource
					obj, err := internal.GetResource(&rcConfig, internal.Kind(kind), id)
					if err != nil {
						return err
					}
					json, err := json.MarshalIndent(obj, "", "  ")
					if err != nil {
						return err
					}
					fmt.Fprintln(c.App.Writer, string(json))
					return nil
				},
			},
			{
				Name:      "apply",
				Usage:     "Apply configurations by filename",
				ArgsUsage: "filename...",
				Action: func(c *cli.Context) error {
					l := c.Args().Len()
					if l == 0 {
						return fmt.Errorf("no files given")
					}
					resources := []internal.GenericResource{}
					for i := 0; i < l; i++ {
						fn := c.Args().Get(i)
						log.Debug().Str("file", fn).Msg("reading file")
						f, err := os.Open(fn)
						if err != nil {
							return err
						}
						bytes, err := ioutil.ReadAll(f)
						if err != nil {
							return err
						}
						log.Trace().RawJSON("job", bytes).Msg("parsing file")
						var parsed internal.GenericResource
						err = json.Unmarshal(bytes, &parsed)
						if err != nil {
							return err
						}
						resources = append(resources, parsed)
					}

					err := internal.Apply(writers, &rcConfig, resources...)
					return err
				},
			},
			{
				Name:  "api-resources",
				Usage: "Print the supported API resources on the server",
				Action: func(c *cli.Context) error {
					kinds, err := internal.GetAllKinds(&rcConfig)
					if err != nil {
						return err
					}
					fmt.Fprintln(c.App.Writer, "NAME")
					for _, kind := range kinds {
						fmt.Fprintln(c.App.Writer, kind)
					}
					return nil
				},
			},
			{
				Name:  "daq",
				Usage: "DAQ related commands",
				Subcommands: []*cli.Command{
					{
						Name:      "autonomous",
						Aliases:   []string{"a"},
						Usage:     "enable or disable autonomous mode on given daq-apps",
						ArgsUsage: "isenabled daq-app-names...",
						Action: func(c *cli.Context) error {
							rawEnabled := c.Args().Get(0)
							enabled := false
							if rawEnabled == "" {
								return fmt.Errorf("bad usage - no isenabled setting")
							}
							if strings.ToLower(rawEnabled) == "false" || strings.ToLower(rawEnabled) == "no" || rawEnabled == "0" {
								enabled = false
							} else if strings.ToLower(rawEnabled) == "true" || strings.ToLower(rawEnabled) == "yes" || rawEnabled == "0" {
								enabled = true
							} else {
								return fmt.Errorf("bad usage - invalid isenabled setting")
							}

							appnames := []string{}
							for i := 1; c.Args().Get(i) != ""; i++ {
								appnames = append(appnames, c.Args().Get(i))
							}
							if len(appnames) == 0 {
								return fmt.Errorf("bad usage - no daq app names given")
							}
							err := daq.SetAutonomousMode(writers, &rcConfig, enabled, appnames...)
							if err != nil {
								return err
							}
							return nil
						},
					},
					{
						Name:  "command",
						Usage: "send command to DAQ Application",
						Flags: []cli.Flag{
							// &cli.StringFlag{Name: "daq-config"},
							&cli.Uint64Flag{Name: "run-number", Usage: "set new run number before running command"},
							&cli.DurationFlag{Name: "timeout", Value: 60 * time.Second, Usage: "how long to wait for the command to complete"},
						},
						ArgsUsage: "command daq-app-names...",
						Action: func(c *cli.Context) error {
							command := c.Args().Get(0)
							if command == "" {
								return fmt.Errorf("bad usage - no command given")
							}
							appnames := []string{}
							for i := 1; c.Args().Get(i) != ""; i++ {
								appnames = append(appnames, c.Args().Get(i))
							}
							if len(appnames) == 0 {
								return fmt.Errorf("bad usage - no daq app names given")
							}
							// overrideConfig := c.String("daq-config")
							// if overrideConfig != "" {
							// 	fmt.Fprintln(c.App.ErrWriter, "warning: custom daq config specified")
							// }
							wg := sync.WaitGroup{}
							wg.Add(len(appnames))
							allSucceeded := true
							for _, name := range appnames {
								go func(name string) {
									log.Debug().Str("name", name).Str("command", command).Msg("preparing to send command")
									args := daq.SendCommandOpts{
										Rc:       &rcConfig,
										DAQApp:   name,
										Command:  command,
										Timeout:  c.Duration("timeout"),
										NewRunNr: nil,
									}
									if c.IsSet("run-number") {
										nr := c.Uint64("run-number")
										args.NewRunNr = &nr
									}
									err := daq.SendCommand(writers, args)
									if err != nil {
										allSucceeded = false
										log.Error().Err(err).Str("name", name).Msg("command failed")
									}
									wg.Done()
								}(name)
							}
							wg.Wait()
							if !allSucceeded {
								return fmt.Errorf("not all DAQ commands succeeded")
							}
							return nil
						},
					},
				},
			},
		},
		Flags:                flags,
		EnableBashCompletion: true,
		HideHelp:             false,
		HideVersion:          false,
		Before: func(c *cli.Context) error {
			err := altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("config-file"))(c)
			if err != nil {
				return err
			}
			l := c.String("log-level")
			if l == "trace" {
				zerolog.SetGlobalLevel(zerolog.TraceLevel)
			} else if l == "debug" {
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			} else if l == "info" {
				zerolog.SetGlobalLevel(zerolog.InfoLevel)
			} else if l == "warn" {
				zerolog.SetGlobalLevel(zerolog.WarnLevel)
			} else if l == "error" {
				zerolog.SetGlobalLevel(zerolog.ErrorLevel)
			} else if l == "fatal" {
				zerolog.SetGlobalLevel(zerolog.FatalLevel)
			} else if l == "panic" {
				zerolog.SetGlobalLevel(zerolog.PanicLevel)
			} else if l != "" {
				log.Warn().Str("requested", l).Str("fallback", zerolog.GlobalLevel().String()).Msg("unknown log level")
			}
			rcConfig.Host = c.String("rc-host")
			rcConfig.ConsulPort = uint32(c.Uint("consul-port"))
			rcConfig.NomadPort = uint32(c.Uint("nomad-port"))

			writers.Out = c.App.Writer
			writers.Err = c.App.ErrWriter

			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(app.ErrWriter, "command failed, %s\n", err)
		os.Exit(1)
	}
}

func wopAction(c *cli.Context) error {
	fmt.Fprintf(c.App.Writer, ":wave: over here, eh\n")
	return nil
}
