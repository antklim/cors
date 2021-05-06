package cors

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	rulesDlm  string = "\n"
	fieldsDlm string = ";"
	valuesDlm string = ","

	fields int = 4

	wildcard string = "*"
)

const (
	pIdx = iota
	oIdx
	hIdx
	mIdx
)

// Methods excluding CONNECT, OPTIONS, TRACE
var allMethods = []string{
	http.MethodDelete,
	http.MethodGet,
	http.MethodHead,
	http.MethodPatch,
	http.MethodPost,
	http.MethodPut,
}

type rule struct {
	o []string // origins
	h []string // headers
	m []string // methods
}

type rules struct {
	raw string
	pr  map[string]rule
}

func newRules(config string) *rules {
	return &rules{raw: config}
}

func (r *rules) Parse() error {
	if strings.TrimSpace(r.raw) == "" {
		return errors.New("invalid cors rules: cannot be empty")
	}

	rawRules := strings.Split(r.raw, rulesDlm)
	r.pr = make(map[string]rule)

	for i, rr := range rawRules {
		pohm := strings.Split(rr, fieldsDlm)

		if s := len(pohm); s != fields {
			return fmt.Errorf("invalid cors rules: invalid amount of fields in rule %d, got %d want %d", i+1, s, fields)
		}

		p := strings.Split(pohm[pIdx], valuesDlm) // paths
		o := strings.Split(pohm[oIdx], valuesDlm) // origins
		h := strings.Split(pohm[hIdx], valuesDlm) // headers

		var m []string // methods
		if pohm[mIdx] == wildcard {
			m = allMethods
		} else {
			m = strings.Split(pohm[mIdx], valuesDlm)
		}

		// TODO: check whether path already registered
		// TODO: break when met wildcard path
		for _, v := range p {
			r.pr[v] = rule{
				o: o,
				h: h,
				m: m,
			}
		}
	}

	return nil
}
