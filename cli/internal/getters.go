package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// GetResource gets the content of a specific resource
func GetResource(c *RCConfig, kind Kind, id string) (*GenericResource, error) {
	if kind == NomadKind {
		return GetNomadJob(c, id)
	}

	url := fmt.Sprintf("http://%s:%d/v1/kv/%ss/%s?raw=", c.Host, c.ConsulPort, kind, id)
	log.Debug().Str("url", url).Str("kind", string(kind)).Str("id", id).Msg("fetching key")
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
		return nil, fmt.Errorf("no such resource exists")
	}

	var parsed GenericResource
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

// GetAllKeys gives all keys from a given kind
func GetAllKeys(c *RCConfig, kind Kind) ([]string, error) {
	if kind == NomadKind {
		return GetAllNomadJobs(c)
	}

	url := fmt.Sprintf("http://%s:%d/v1/kv/%ss/?keys=&separator=/", c.Host, c.ConsulPort, kind)
	log.Debug().Str("url", url).Str("kind", string(kind)).Msg("fetching all keys")
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
		return []string{}, nil
	}

	var parsed []string
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		return nil, err
	}

	result := []string{}
	for _, p := range parsed {
		if !strings.HasSuffix(p, "s/") {
			result = append(result, p[len(kind)+2:])
		}
	}
	return result, nil
}

// GetAllKinds gives all known and found kinds
func GetAllKinds(c *RCConfig) ([]Kind, error) {
	result := []Kind{DAQAppKind, NomadKind}

	url := fmt.Sprintf("http://%s:%d/v1/kv/?keys=&separator=/", c.Host, c.ConsulPort)
	log.Debug().Str("url", url).Msg("fetching all kinds")
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

	var parsed []string
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		return nil, err
	}

	for _, p := range parsed {
		if strings.HasSuffix(p, "s/") {
			p = p[:len(p)-2]
			alreadyIn := false
			for _, r := range result {
				if string(r) == p {
					alreadyIn = true
					break
				}
			}
			if !alreadyIn {
				result = append(result, Kind(p))
			}
		}
	}

	return result, nil
}
