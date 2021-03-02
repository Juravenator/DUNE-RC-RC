package internal

import "io"

// Writers to write to CLI output
type Writers struct {
	Out io.Writer
	Err io.Writer
}

// GenericResource is the generic structure all resources conform to
type GenericResource struct {
	Meta   GenericMeta `json:"meta"`
	Spec   GenericSpec `json:"spec"`
	Status interface{} `json:"status,omitempty"`
}

// GenericMeta is part of GenericResource
type GenericMeta struct {
	Kind Kind   `json:"kind"`
	Name string `json:"name"`
}

// GenericSpec is part of GenericResource
type GenericSpec map[string]interface{}

// Kind hardcodes known resource kinds
type Kind string

// known kinds
const (
	NomadKind     Kind = "nomad-job"
	DAQAppKind    Kind = "daq-application"
	DAQConfigKind Kind = "daq-config"
)
