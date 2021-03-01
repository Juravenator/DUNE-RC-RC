package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"cli.rc.ccm.dunescience.org/internal"

	"github.com/rs/zerolog"
	logger "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

var log = logger.With().Str("pkg", "internal").Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})

var rcConfig = internal.RCConfig{}

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
						for _, kind := range kinds {
							fmt.Fprintln(c.App.Writer, "KEY")
							keys, err := internal.GetAllKeys(&rcConfig, kind)
							if err != nil {
								return err
							}
							for _, key := range keys {
								fmt.Fprintln(c.App.Writer, kind+"s/"+key)
							}
							fmt.Fprintln(c.App.Writer, "")
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
					obj, err := internal.GetResource(&rcConfig, kind, id)
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