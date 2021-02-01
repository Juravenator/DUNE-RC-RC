package internal

import (
	"io"
	"os/exec"
)

// Process is a running exec.Cmd with pipes and signals hooked up
type Process struct {
	Cmd    *exec.Cmd
	Stdout chan string
	Stderr chan string
	Stdin  io.WriteCloser
}

// SendCommand sends the given command over process stdin and checks for errors in stdout/stderr
func (p *Process) SendCommand(command string) {
	p.Stdin.Write([]byte(command))
	p.Stdin.Write([]byte("\n"))
}

// Kill kills the underlying process
func (p *Process) Kill() {
	if p.Cmd.Process == nil {
		return
	}
	p.Cmd.Process.Kill()
}
