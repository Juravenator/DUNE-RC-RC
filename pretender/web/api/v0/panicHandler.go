package apiv0

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"

	"dune.rc.pretender/web/api/v0/util"

	"github.com/go-chi/chi/middleware"
)

func panicHandler(next http.Handler) http.Handler {
	handlePanic := func(w http.ResponseWriter, r *http.Request) {
		// https://blog.golang.org/defer-panic-and-recover
		defer panicRecoverer(w, r)
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(handlePanic)
}

func panicRecoverer(w http.ResponseWriter, r *http.Request) {
	if panicContent := recover(); panicContent != nil && panicContent != http.ErrAbortHandler {
		logEntry := middleware.GetLogEntry(r)
		if logEntry != nil {
			logEntry.Panic(panicContent, debug.Stack())
		} else {
			fmt.Fprintf(os.Stderr, "Panic: %+v\n", panicContent)
			debug.PrintStack()
		}

		util.Error(http.StatusInternalServerError).Send(w, r)
	}
}
