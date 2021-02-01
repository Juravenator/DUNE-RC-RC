package api

import (
	"net"
	"net/http"

	"dune.rc.pretender/internal"
	v0 "dune.rc.pretender/web/api/v0"

	"github.com/go-chi/chi"
)

// Setup opens a socket on the specified port, but does not accept requests yet
func Setup(port string, process *internal.DAQProcess) (net.Listener, chi.Router, error) {
	r := chi.NewRouter()

	v0.Register(r, process)

	if port == "" {
		port = "0"
	}
	listener, err := net.Listen("tcp", ":"+port)
	return listener, r, err
}

// Start starts the API
func Start(listener net.Listener, r chi.Router) error {
	address := "http://" + listener.Addr().String()
	log.WithField("address", address).Info("listening")
	return http.Serve(listener, r)
}
