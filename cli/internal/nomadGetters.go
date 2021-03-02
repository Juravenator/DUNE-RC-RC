package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// NomadWrapperStatus is part of NomadWrapper
type NomadWrapperStatus struct {
	Status  string  `json:"status"`
	Stable  bool    `json:"stable"`
	Version float64 `json:"version"`
}

// GetAllNomadJobs gives all names of Nomad jobs
func GetAllNomadJobs(c *RCConfig) ([]string, error) {
	url := fmt.Sprintf("http://%s:%d/v1/jobs", c.Host, c.NomadPort)
	log.Debug().Str("url", url).Msg("fetching all keys")
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Trace().Bytes("body", body).Msg("received response")
	if len(body) == 0 {
		return nil, fmt.Errorf("no such kind exists")
	}

	var parsed []GenericSpec
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		return nil, err
	}

	result := []string{}
	for _, p := range parsed {
		result = append(result, p["ID"].(string))
	}
	return result, nil
}

// GetNomadJob gives a given jobs configuration
func GetNomadJob(c *RCConfig, id string) (*GenericResource, error) {
	url := fmt.Sprintf("http://%s:%d/v1/job/%s", c.Host, c.NomadPort, id)
	log.Debug().Str("url", url).Str("id", id).Msg("fetching all keys")
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Trace().Bytes("body", body).Msg("received response")
	if string(body) == "job not found" {
		return nil, fmt.Errorf("no such job exists")
	}

	var parsed GenericSpec
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		return nil, err
	}

	result := GenericResource{
		Meta: GenericMeta{
			Kind: NomadKind,
			Name: id,
		},
		Spec: parsed,
		Status: NomadWrapperStatus{
			Status:  parsed["Status"].(string),
			Stable:  parsed["Stable"].(bool),
			Version: parsed["Version"].(float64),
		},
	}
	return &result, nil
}
