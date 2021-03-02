package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type namedHandler struct {
	name   string
	server *http.Server
}

func (h namedHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	defer func() {
		go h.server.Shutdown(context.Background())
	}()
	if strings.ToUpper(req.Method) != "POST" {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(resp, "only POST allowed")
		return
	}
	defer req.Body.Close()
	bytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error().Err(err).Msg("reading return path response failed")
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	var parsed struct {
		Result string `json:"result"`
	}
	err = json.Unmarshal(bytes, &parsed)
	if err != nil {
		log.Error().Err(err).Bytes("body", bytes).Msg("parsing return path response failed")
		resp.WriteHeader(http.StatusBadRequest)
	}
	log.Error().Err(err).Bytes("body", bytes).Str("result", parsed.Result).Msg("response received")
	fmt.Fprintf(resp, "this is endpoint %s: %t", h.name, parsed.Result == "OK")
}

func randomListener(name string) error {
	listener, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		return err
	}
	fmt.Printf("server %s listening on %s\n", name, listener.Addr().String())

	server := http.Server{}
	handler := namedHandler{name, &server}
	server.Handler = handler
	server.SetKeepAlivesEnabled(false)

	serverErr := make(chan error)
	go func() { serverErr <- server.Serve(listener) }()
	select {
	case err := <-serverErr:
		return err
	case <-time.After(10 * time.Second):
		server.Close()
		return fmt.Errorf("too late for %s", name)
	}
}

func main() {
	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		err := randomListener("two")
		fmt.Printf("%s exited with %s\n", "two", err)
		wg.Done()
	}()
	go func() {
		err := randomListener("two")
		fmt.Printf("%s exited with %s\n", "two", err)
		wg.Done()
	}()

	wg.Wait()
	fmt.Println("finished!")
}
