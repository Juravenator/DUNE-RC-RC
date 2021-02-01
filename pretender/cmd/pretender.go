package main

import (
	"dune.rc.pretender/internal"
	"dune.rc.pretender/web/api"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("starting")

	/* get CLI args */
	args, err := internal.ParseCLI()
	if err != nil {
		log.WithError(err).Fatal("cannot parse arguments")
	}
	log.WithField("args", args).Info("parsed arguments")

	// we have two blocking actions that can fail
	// API & subprocess
	errChan := make(chan error, 2)

	/* start wrapped process */
	process, err := internal.StartWrapper(args.WrappedCommand, args.WrappedArgs)
	if err != nil {
		log.WithError(err).Fatal("cannot start subprocess")
	}
	defer process.Kill()
	log.Info("subprocess started")
	go func() {
		errChan <- process.Cmd.Wait()
		close(process.Stdout)
		close(process.Stderr)
	}()
	daq := internal.ToDAQProcess(process)

	/* start API */
	listener, router, err := api.Setup(args.Port, &daq)
	if err != nil {
		log.WithError(err).Fatal("cannot setup API")
	}
	defer listener.Close()
	log.Info("starting API")
	go func() {
		errChan <- api.Start(listener, router)
	}()

	/* wait for something to fail, or subprocess to exit cleanly */
	err = <-errChan
	if err == nil {
		log.Info("subprocess exited without error")
	} else {
		log.WithError(err).Error("error received, shutting down")
	}

	log.Info("stopping wrapper process")
}
