package internal

import "fmt"

type ServiceAddress struct {
	Address string
	Port    uint16
}

func GetFirstServiceAddress(c *RCConfig, serviceName string) (*ServiceAddress, error) {
	all, err := GetServiceAddresses(c, serviceName)
	if err != nil {
		return nil, err
	}
	if len(all) == 0 {
		return nil, fmt.Errorf("no such service")
	}
	return &(all[0]), nil
}
