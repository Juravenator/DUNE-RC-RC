package apiv0

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"dune.rc.pretender/internal"
	"dune.rc.pretender/web/api/v0/state"
	"dune.rc.pretender/web/api/v0/util"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

var v0router = makeRouter()

// Register adds the v0router to the given router
func Register(router *chi.Mux, process *internal.DAQProcess) {
	v0router.Get("/", listAllEndpoints)
	router.Mount("/api/v0", v0router)
	state.Register(v0router, process)
}

func makeRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(panicHandler)
	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped
	r.Use(middleware.Timeout(5 * time.Minute))

	r.NotFound(statusReturner(http.StatusNotFound))
	r.MethodNotAllowed(statusReturner(http.StatusMethodNotAllowed))
	return r
}

func statusReturner(code int) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		util.Error(code).Send(w, r)
	}
}

func listAllEndpoints(w http.ResponseWriter, r *http.Request) {
	endpointsMap := make(map[string][]string)

	walker := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		entry := endpointsMap[route]
		if entry == nil {
			entry = []string{method}
		} else {
			entry = append(entry, method)
		}
		endpointsMap[route] = entry
		return nil
	}
	err := chi.Walk(v0router, walker)
	if err != nil {
		util.Error(500).Send(w, r)
		return
	}

	endpoints := make([]string, len(endpointsMap))

	i := 0
	for route, methods := range endpointsMap {
		sort.Strings(methods)
		endpoints[i] = fmt.Sprintf("(%s) %s", strings.Join(methods, ","), route)
		i++
	}

	sort.Strings(endpoints)

	render.JSON(w, r, endpoints)
}
