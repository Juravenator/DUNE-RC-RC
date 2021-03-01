package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// GetResource gets the content of a specific resource
func GetResource(c *RCConfig, kind string, id string) (interface{}, error) {
	url := fmt.Sprintf("http://%s:%d/v1/kv/%ss/%s?raw=", c.Host, c.ConsulPort, kind, id)
	log.Debug().Str("url", url).Str("kind", kind).Str("id", id).Msg("fetching key")
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var parsed interface{}
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		return nil, err
	}
	return parsed, nil
}

// GetAllKeys gives all keys from a given kind
func GetAllKeys(c *RCConfig, kind string) ([]string, error) {
	url := fmt.Sprintf("http://%s:%d/v1/kv/%ss/?keys=&separator=/", c.Host, c.ConsulPort, kind)
	log.Debug().Str("url", url).Str("kind", kind).Msg("fetching all keys")
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
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
func GetAllKinds(c *RCConfig) ([]string, error) {
	result := []string{"daq-application"}

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
				if r == p {
					alreadyIn = true
					break
				}
			}
			if !alreadyIn {
				result = append(result, p)
			}
		}
	}

	return result, nil
}
