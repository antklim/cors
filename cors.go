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
	for _, p := range paths {
		h := handlers.CORS(
			handlers.AllowedHeaders([]string{"*"}),
			handlers.AllowedMethods([]string{"DELETE"}),
		)(noopHTTPHandler)
		router.Handle(p, h).Methods(http.MethodOptions)
	}

	return router, nil
}
