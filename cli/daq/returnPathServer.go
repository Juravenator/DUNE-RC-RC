package daq

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

// ReturnPathHandler wraps handler for DAQ return-path specific needs
type ReturnPathHandler struct {
	server          *http.Server
	DAQResponseChan chan string
}

func (h ReturnPathHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	log.Debug().Str("method", req.Method).Str("url", req.URL.String()).Msg("received request")
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
	log.Debug().Err(err).Bytes("body", bytes).Str("result", parsed.Result).Msg("response received")

	resp.WriteHeader(http.StatusOK)
	resp.Header().Set("Content-type", "text/plain")
	fmt.Fprintf(resp, "Response received")

	h.DAQResponseChan <- parsed.Result
}

func openReturnPath(l net.Listener, timeout time.Duration) (*string, error) {
	server := http.Server{}
	handler := ReturnPathHandler{&server, make(chan string)}
	server.Handler = handler
	server.SetKeepAlivesEnabled(false)

	go server.Serve(l)
	select {
	case daqResponse := <-handler.DAQResponseChan:
		return &daqResponse, nil
	case <-time.After(timeout):
		server.Close()
		return nil, fmt.Errorf("return path listener timed out")
	}
}
