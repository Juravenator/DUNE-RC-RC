package state

import (
	"encoding/json"
	"errors"
	"net/http"

	"dune.rc.pretender/internal"
	"dune.rc.pretender/web/api/v0/util"
	"github.com/go-chi/render"
)

func init() {
	router.Get("/", getState)
	router.Post("/", postCommand)
}

type stateResponse struct {
	State         internal.DAQState `json:"state"`
	Transitioning bool              `json:"transitioning"`
}

func getState(w http.ResponseWriter, r *http.Request) {
	response := stateResponse{
		State:         daq.State,
		Transitioning: daq.TransitioningState,
	}
	render.JSON(w, r, response)
}

type commandRequest struct {
	Command string                 `json:"command"`
	Payload map[string]interface{} `json:"payload,omitempty"`
}

func postCommand(w http.ResponseWriter, r *http.Request) {
	var payload commandRequest
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		util.ErrorMsg(400, "cannot parse request body")
	}
	err = daq.SendCommandAndWait(payload.Command)
	if err != nil {
		code := 500
		if errors.Is(err, internal.CannotStateChangeError) {
			code = 400
		}
		util.ErrorMsg(code, err.Error()).Send(w, r)
		return
	}
	util.Reply(202).Send(w, r)
}
