package internal

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type ServiceResponse struct {
	Service ServiceResponseService
}

type ServiceResponseService struct {
	Address string
	Port int
}

func GetServiceAddresses(serviceName string) []ServiceAddress {
	resp, err := http.Get("http://localhost:8500/v1/health/service/" + serviceName + "?dc=dune-rc")
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
			panic(err.Error())
	}

	var parsed []ServiceResponse
	err = json.Unmarshal(body, &parsed)
	if err != nil {
			panic(err.Error())
	}

	var result = []ServiceAddress{}
	for _, s := range parsed {
		result = append(result, ServiceAddress{Address: s.Service.Address, Port: s.Service.Port})
	}
	return result
}