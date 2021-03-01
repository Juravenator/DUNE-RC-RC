package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// NomadJob (very lazily) describes a Nomad Job object
type NomadJob map[string]interface{}

// NomadWrapper wraps around a bare nomad job for consistency
type NomadWrapper struct {
	Meta   NomadWrapperMeta   `json:"meta"`
	Spec   NomadJob           `json:"spec"`
	Status NomadWrapperStatus `json:"status,omitempty"`
}

// NomadKindType to hardcode kind
type NomadKindType string

// NomadKind hardcodes nomad kind
const NomadKind NomadKindType = "nomad-job"

// NomadWrapperMeta is part of NomadWrapper
type NomadWrapperMeta struct {
	Kind NomadKindType `json:"kind"`
	Name string        `json:"name"`
}

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

	var parsed []NomadJob
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
func GetNomadJob(c *RCConfig, id string) (*NomadWrapper, error) {
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

	var parsed NomadJob
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		return nil, err
	}

	result := NomadWrapper{
		Meta: NomadWrapperMeta{
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
