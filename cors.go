package cors

import (
	"errors"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// TODO: support different rules config format: yaml, json

// Routes(<list of paths>, <cors config>) http.Handler
//
// Cors config format: ruleA\nruleB...\nruleX
//
// Rule format: PATHs;ORIGINs;HEADERs;METHODs
// path can be *
// allowed origins can be *
// allowed headers should be explicit
// allowed methods can be *

var noopHTTPHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

func OptionsRoutes(paths []string, config string) (http.Handler, error) {
	if len(paths) == 0 {
		return nil, errors.New("invalid paths list: cannot be empty")
	}

	r := NewRules(config)
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
			// TODO: delete rule here to avoid duplication when reached a wildcard
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

func addRule(router *mux.Router, path string, r Rule) {
	h := handlers.CORS(
		handlers.AllowedHeaders(r.h),
		handlers.AllowedMethods(r.m),
		handlers.AllowedOrigins(r.o),
	)(noopHTTPHandler)
	router.Handle(path, h).Methods(http.MethodOptions)
}

// TODO: implement
// func RouteMiddleware(path string, r Rule) func(http.Handler) http.Handler {
// }
