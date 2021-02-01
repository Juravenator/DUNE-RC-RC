package internal

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"

	log "github.com/sirupsen/logrus"
)

// StartWrapper sets up and runs the command to wrap
func StartWrapper(command string, args []string) (*Process, error) {
	cmd := exec.Command(command, args...)
	process, err := setupProcess(cmd)
	if err != nil {
		return nil, err
	}
	return process, nil
}

// hook stdin, stdout, stderr, and forward signals
func setupProcess(cmd *exec.Cmd) (*Process, error) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("cannot attach to stdout: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("cannot attach to stderr: %w", err)
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("cannot attach to stdin: %w", err)
	}

	outChan := make(chan string, 10)
	errChan := make(chan string, 10)

	go linesToChannel(stdout, outChan)
	go linesToChannel(stderr, errChan)

	// forwardSignals relies on the process already being started
	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("cannot start process: %w", err)
	}

	go forwardSignals(cmd)

	return &Process{Cmd: cmd, Stdout: outChan, Stderr: errChan, Stdin: stdin}, nil
}

func linesToChannel(input io.Reader, ch chan string) {
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		select {
		case ch <- line: // Put line in the channel unless it is full
		default:
			log.Warn("channel full, discarding")
		}
		os.Stdout.WriteString(line)
		os.Stdout.WriteString("\n")
	}
}

func forwardSignals(cmd *exec.Cmd) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan)
	for {
		sig := <-sigChan
		if sig.String() == "child exited" {
			continue
		}
		log.WithField("sig", sig).Info("forwarding signal to subprocess")
		err := cmd.Process.Signal(sig)
		if err != nil {
			log.WithField("sig", sig).Error("failed to signal subprocess")
		}

	}
}
