package internal

import (
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var CannotStateChangeError = fmt.Errorf("cannot change state")
var timeoutError = fmt.Errorf("time out while checking log output")

// DAQState is a known DAQ state
type DAQState string

// Known DAQStates
// stolen from logs from running https://github.com/DUNE-DAQ/minidaqapp/wiki/Simple-instructions-for-running-the-app
// Available commands: | init | conf | start | stop | pause | resume | scrap
const (
	UnknownState    DAQState = "UNKNOWN"
	ErrorState      DAQState = "ERROR"
	InitState       DAQState = "INIT"
	ConfiguredState DAQState = "CONFIGURED"
	StartedState    DAQState = "STARTED"
	PausedState     DAQState = "PAUSED"
	StoppedState    DAQState = "STOPPED"
	ScrappedState   DAQState = "SCRAPPED"
)

// DAQProcess wraps Process with DAQ state
type DAQProcess struct {
	*Process
	State              DAQState
	TransitioningState bool
}

// ToDAQProcess turns a Process into a DAQProcess
func ToDAQProcess(process *Process) DAQProcess {
	return DAQProcess{process, UnknownState, false}
}

// ChangeState issues commands to push the DAQ process in a certain state
// it does NOT block for the command to succeeded, check p.TransitioningState and p.State
// func (p *DAQProcess) ChangeState(desiredState DAQState) error {

// 	// figure out what command we want to send
// 	command := ""
// 	switch desiredState {
// 	case InitState:
// 		command = "init"
// 	case ConfiguredState:
// 		command = "conf"
// 	case StartedState:
// 		command = "start"
// 		if p.State == PausedState {
// 			command = "resume"
// 		}
// 	case PausedState:
// 		command = "pause"
// 	case StoppedState:
// 		command = "stop"
// 	case ScrappedState:
// 		command = "scrap"
// 	}

// 	if command == "" {
// 		return fmt.Errorf("Cannot determine command to change to state '%s': %w", desiredState, CannotStateChangeError)
// 	}
func (p *DAQProcess) SendCommandAndWait(command string) error {
	var desiredStateFromCommand DAQState
	switch command {
	case "init":
		desiredStateFromCommand = InitState
	case "conf":
		desiredStateFromCommand = ConfiguredState
	case "start", "resume":
		desiredStateFromCommand = StartedState
	case "pause":
		desiredStateFromCommand = PausedState
	case "stop":
		desiredStateFromCommand = StoppedState
	case "scrap":
		desiredStateFromCommand = ScrappedState
	default:
		return fmt.Errorf("Cannot determine desired state from command '%s': %w", command, CannotStateChangeError)
	}

	// drain any unread log messages
	drain(p.Stdout)
	drain(p.Stderr)

	// engage
	p.TransitioningState = true
	p.SendCommand(command)

	// in background, check logs to determine success or failure
	go func() {
		done := make(chan bool, 1)
		badOut := make(chan error, 1)
		go checkForLogStateChange(p.Process.Stdout, badOut, done)
		go checkForLogStateChange(p.Process.Stderr, badOut, done)
		for err := range badOut {
			if !p.TransitioningState { // draining
				continue
			}
			if err == nil {
				p.State = DAQState(desiredStateFromCommand)
			} else {
				p.State = ErrorState
				log.WithError(err).Error("possible command failure")
			}
			p.TransitioningState = false
			break
		}
		close(badOut)
		close(done)

	}()
	return nil
}

// watch logs until feedback about given command is found in the given output stream
func checkForLogStateChange(input chan string, badOut chan error, done chan bool) {
	for {
		select {
		case <-done:
			return
		case line := <-input:
			if strings.Contains(line, "Command execution resulted with") {
				if !strings.Contains(line, "Command execution resulted with: OK") {
					badOut <- fmt.Errorf("Error in logs: %s", line)
				} else {
					badOut <- nil
				}
				done <- true
				return
			}
			// pieces := strings.SplitN(line, " ", 4)
			// if len(pieces) == 4 && pieces[2] == "ERROR" {
			// 	badOut <- fmt.Errorf("Error in logs: %s", line)
			// }
		case <-time.After(30 * time.Second):
			badOut <- timeoutError
			return
		}
	}
}

func drain(c chan string) {
	log.Debug("draining")
L:
	for {
		select {
		case <-c:
		default:
			break L
		}
	}
	log.Debug("draining finished")
}
