package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ServiceResponse struct {
	Service ServiceResponseService
}

type ServiceResponseService struct {
	Address string
	Port    uint16
}

func GetServiceAddresses(c *RCConfig, serviceName string) ([]ServiceAddress, error) {
	url := fmt.Sprintf("http://%s:%d/v1/health/service/%s?dc=dune-rc", c.Host, c.ConsulPort, serviceName)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var parsed []ServiceResponse
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		return nil, err
	}

	var result = []ServiceAddress{}
	for _, s := range parsed {
		result = append(result, ServiceAddress{Address: s.Service.Address, Port: s.Service.Port})
	}
	return result, nil
}
