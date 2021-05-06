package cors

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type rule struct {
	raw string
	p   []string // paths
	o   []string // origins
	h   []string // headers
	m   []string // methods
}

type rules struct {
	raw string
	r   []rule
}

func newRules(config string) rules {
	return rules{raw: config}
}

func (r rules) Validate() error {
	if strings.TrimSpace(r.raw) == "" {
		return errors.New("invalid cors rules: cannot be empty")
	}
	return nil
}

var noopHTTPHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

func Routes(paths []string, config string) (http.Handler, error) {
	if len(paths) == 0 {
		return nil, errors.New("invalid paths list: cannot be empty")
	}

	r := newRules(config)
	if err := r.Validate(); err != nil {
		return nil, err
	}

	router := mux.NewRouter()
	for _, p := range paths {
		h := handlers.CORS()(noopHTTPHandler)
		router.Handle(p, h).Methods(http.MethodOptions)
	}

	return router, nil
}
