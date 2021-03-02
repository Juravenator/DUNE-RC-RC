package internal

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Apply given files by name
func Apply(w Writers, c *RCConfig, resources ...GenericResource) error {
	log.Debug().Msg("applying resources")
	consulTodo := []GenericResource{}
	nomadTodo := []GenericSpec{}
	for _, r := range resources {
		if r.Meta.Kind == NomadKind {
			if r.Meta.Name != r.Spec["ID"] || r.Spec["ID"] != r.Spec["Name"] {
				return fmt.Errorf("meta.name should equal spec.Name and spec.ID, got '%s', '%s', '%s'", r.Meta.Name, r.Spec["Name"], r.Spec["ID"])
			}
			nomadTodo = append(nomadTodo, r.Spec)
		} else {
			consulTodo = append(consulTodo, r)
		}
	}
	log.Debug().Int("consuljobs", len(consulTodo)).Int("nomadjobs", len(nomadTodo)).Msg("parsed files")

	// consul jobs go first (nomad configs might depend on them)
	fmt.Fprint(w.Out, "committing consul transaction... ")
	err := consulTransaction(c, consulTodo...)
	if err != nil {
		return err
	}
	fmt.Fprintln(w.Out, "OK")

	// nomad jobs go after
	err = runNomad(w, c, nomadTodo...)
	if err != nil {
		return err
	}

	return nil
}

// ConsulTransaction contains a key-value transaction
type ConsulTransaction struct {
	KV ConsulKVTransaction `json:"KV"`
}

// ConsulKVTransaction is part of ConsulTransaction
type ConsulKVTransaction struct {
	Verb  string `json:"Verb"`
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

func consulTransaction(c *RCConfig, stuff ...GenericResource) error {
	transactions := []ConsulTransaction{}
	for _, s := range stuff {
		b, err := json.Marshal(s)
		if err != nil {
			return err
		}
		t := ConsulKVTransaction{
			Verb:  "set",
			Key:   fmt.Sprintf("%ss/%s", s.Meta.Kind, s.Meta.Name),
			Value: base64.StdEncoding.EncodeToString(b),
		}
		transactions = append(transactions, ConsulTransaction{t})
	}

	payload, err := json.Marshal(transactions)
	if err != nil {
		return err
	}
	client := &http.Client{}
	url := fmt.Sprintf("http://%s:%d/v1/txn", c.Host, c.ConsulPort)
	log.Trace().Bytes("req", payload).Str("url", url).Msg("registering nomad job")
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("consul transaction failed with code %s", resp.Status)
	}
	return nil
}

// NomadRequest describes a request body to nomad
type NomadRequest struct {
	Job GenericSpec `json:"Job"`
}

func runNomad(w Writers, c *RCConfig, jobs ...GenericSpec) error {
	for _, job := range jobs {
		fmt.Fprintf(w.Out, "running nomad job %s... ", job["Name"])
		b, err := json.Marshal(NomadRequest{job})
		if err != nil {
			return err
		}
		url := fmt.Sprintf("http://%s:%d/v1/job/%s", c.Host, c.NomadPort, job["Name"])
		log.Trace().Bytes("req", b).Str("url", url).Msg("registering nomad job")
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(b))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		log.Trace().Bytes("body", body).Msg("received response")
		if resp.StatusCode != 200 {
			return fmt.Errorf("nomad job failed with code %s", resp.Status)
		}
		fmt.Fprintln(w.Out, "OK")
	}
	return nil
}
