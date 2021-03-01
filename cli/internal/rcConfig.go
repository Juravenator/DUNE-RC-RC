package internal

// RCConfig contains data needed to connect to RC
type RCConfig struct {
	Host       string
	ConsulPort uint32
	NomadPort  uint32
}
