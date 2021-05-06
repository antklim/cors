package cors

import (
	"errors"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var noopHTTPHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

func Routes(paths []string, config string) (http.Handler, error) {
	if len(paths) == 0 {
		return nil, errors.New("invalid paths list: cannot be empty")
	}

	r := newRules(config)
	if err := r.Parse(); err != nil {
		return nil, err
	}

	router := mux.NewRouter()

	for p, rule := range r.pr {
		if p == wildcard {
			for _, pp := range paths {
				addRule(router, pp, rule)
			}
			break
		}

		found, _ := find(paths, p)
		if found {
			addRule(router, p, rule)
		}
	}

	return router, nil
}

func find(a []string, x string) (bool, int) {
	for i, v := range a {
		if v == x {
			return true, i
		}
	}
	return false, 0
}

// func deleteAt(a []string, idx int) []string {
// 	return append(a[:idx], a[idx+1:]...)
// }

func addRule(router *mux.Router, path string, r rule) {
	h := handlers.CORS(
		handlers.AllowedHeaders(r.h),
		handlers.AllowedMethods(r.m),
		handlers.AllowedOrigins(r.o),
	)(noopHTTPHandler)
	router.Handle(path, h).Methods(http.MethodOptions)
}
