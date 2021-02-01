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
	router.Put("/", putState)
}

type stateResponse struct {
	State internal.DAQState `json:"state"`
}

func getState(w http.ResponseWriter, r *http.Request) {
	response := stateResponse{
		State: daq.State,
	}
	render.JSON(w, r, response)
}

type putStateRequest struct {
	State   string                 `json:"state"`
	Payload map[string]interface{} `json:"payload,omitempty"`
}

func putState(w http.ResponseWriter, r *http.Request) {
	var payload putStateRequest
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		util.ErrorMsg(400, "cannot parse request body")
	}
	err = daq.ChangeState(internal.DAQState(payload.State))
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
