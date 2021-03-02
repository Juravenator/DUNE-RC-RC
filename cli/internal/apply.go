package internal

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// GenericResource is the generic structure all resources conform to
type GenericResource struct {
	Meta GenericMeta `json:"meta"`
	Spec GenericSpec `json:"spec"`
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
	NomadKind Kind = "nomad-job"
)

// Apply given files by name
func Apply(writer io.Writer, c *RCConfig, filenames ...string) error {
	log.Debug().Strs("files", filenames).Msg("applying files")
	consulTodo := []GenericResource{}
	nomadTodo := []GenericSpec{}
	for _, n := range filenames {
		f, err := os.Open(n)
		if err != nil {
			return err
		}
		bytes, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}
		log.Trace().RawJSON("job", bytes).Msg("parsing file")
		var parsed GenericResource
		err = json.Unmarshal(bytes, &parsed)
		if err != nil {
			return err
		}
		if parsed.Meta.Kind == NomadKind {
			if parsed.Meta.Name != parsed.Spec["ID"] || parsed.Spec["ID"] != parsed.Spec["Name"] {
				return fmt.Errorf("meta.name should equal spec.Name and spec.ID, got '%s', '%s', '%s'", parsed.Meta.Name, parsed.Spec["Name"], parsed.Spec["ID"])
			}
			nomadTodo = append(nomadTodo, parsed.Spec)
		} else {
			consulTodo = append(consulTodo, parsed)
		}
	}
	log.Debug().Int("consuljobs", len(consulTodo)).Int("nomadjobs", len(nomadTodo)).Msg("parsed files")

	// consul jobs go first (nomad configs might depend on them)
	fmt.Fprint(writer, "committing consul transaction... ")
	err := consulTransaction(c, consulTodo...)
	if err != nil {
		return err
	}
	fmt.Fprintln(writer, "OK")

	// nomad jobs go after
	err = runNomad(writer, c, nomadTodo...)
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

func runNomad(writer io.Writer, c *RCConfig, jobs ...GenericSpec) error {
	for _, job := range jobs {
		fmt.Fprintf(writer, "running nomad job %s... ", job["Name"])
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
		fmt.Fprintln(writer, "OK")
	}
	return nil
}
