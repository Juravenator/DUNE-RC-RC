package state

import (
	"dune.rc.pretender/internal"
	"github.com/go-chi/chi"
)

var router = chi.NewRouter()
var daq *internal.DAQProcess

// Register adds this router to the given parent router
func Register(r *chi.Mux, process *internal.DAQProcess) {
	daq = process
	r.Mount("/state", router)
}
