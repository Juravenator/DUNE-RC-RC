package internal

import (
	"fmt"
	"os"
)

// CLIArgs contains parsed CLI arguments
type CLIArgs struct {
	Port           string
	WrappedCommand string
	WrappedArgs    []string
}

// ParseCLI generates parsed CLI arguments
func ParseCLI() (*CLIArgs, error) {
	a := CLIArgs{
		Port:           "0",
		WrappedCommand: "",
		WrappedArgs:    []string{},
	}
	// parse all own arguments, first argument is our own binary
	i := 1
loop:
	for ; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--":
			break loop
		case "-p":
			i++
			if i == len(os.Args) {
				return nil, fmt.Errorf("no port given after port flag")
			}
			a.Port = os.Args[i]
		default:
			break loop
		}
	}
	// remaining args are wrapper command
	if i == len(os.Args) {
		return nil, fmt.Errorf("no wrapper command given")
	}
	a.WrappedCommand = os.Args[i]
	i++
	// if there are remaining arguments, add them
	if i != len(os.Args) {
		a.WrappedArgs = os.Args[i:]
	}
	return &a, nil
}
