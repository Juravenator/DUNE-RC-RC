// https://github.com/go-chi/chi/blob/baf4ef5b139e284b297573d89daf587457153aa3/_examples/rest/main.go#L399

package util

import (
	"net/http"

	"github.com/go-chi/render"
)

// Response renderer type for JSON API responses for all non-200 HTTP Codes
type Response struct {
	HTTPCode int
	Err      string `json:"error,omitempty"`
	AppErr   string `json:"appError,omitempty"`
	AppCode  int    `json:",omitempty"`
	Msg      string `json:"message,omitempty"`
}

// Send error to a standard http.ResponseWriter
func (e *Response) Send(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	render.Status(r, e.HTTPCode)
	render.JSON(w, r, e)
	return nil
}

// Error constructs a default Response for the given http code
func Error(httpCode int) *Response {
	return ErrorMsg(httpCode, "")
}

// ErrorMsg constructs an Response with the given error
func ErrorMsg(httpCode int, err string) *Response {
	if err == "" {
		err = http.StatusText(httpCode)
	}
	return &Response{
		HTTPCode: httpCode,
		Err:      err,
	}
}

// Reply constructs an Response with the given message
func Reply(httpCode int) *Response {
	return ReplyMsg(httpCode, "")
}

// ReplyMsg constructs an Response with the given message
func ReplyMsg(httpCode int, message string) *Response {
	if message == "" {
		message = http.StatusText(httpCode)
	}
	return &Response{
		HTTPCode: httpCode,
		Msg:      message,
	}
}

// AppInfo adds app-specific error info
func (e *Response) AppInfo(err string, code int) *Response {
	e.AppErr = err
	e.AppCode = code
	return e
}

// UserMessage adds the given message for the end user
func (e *Response) UserMessage(msg string) *Response {
	e.Msg = msg
	return e
}
