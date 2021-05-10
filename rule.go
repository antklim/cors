package cors

import (
	"fmt"
	"net/http"
	"strings"
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

const (
	rulesDlm  string = "\n"
	fieldsDlm string = ";"
	valuesDlm string = ","

	wildcard string = "*"
)

const (
	pIdx = iota
	oIdx
	hIdx
	mIdx
	fNum // a number of fields in the rule, must be last
)

const parseErr string = "invalid cors rules"

// Methods excluding CONNECT, OPTIONS, TRACE
var allMethods = []string{
	http.MethodDelete,
	http.MethodGet,
	http.MethodHead,
	http.MethodPatch,
	http.MethodPost,
	http.MethodPut,
}

var validMethods = append([]string{http.MethodConnect, http.MethodOptions, http.MethodTrace}, allMethods...)

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
		return fmt.Errorf("%s: cannot be empty", parseErr)
	}

	rawRules := strings.Split(r.raw, rulesDlm)
	r.pr = make(map[string]rule)

	for i, rr := range rawRules {
		rr = strings.TrimSpace(rr)

		// skip empty rows in confing
		if len(rr) == 0 {
			continue
		}

		pohm := strings.Split(rr, fieldsDlm)

		if s := len(pohm); s != fNum {
			return fmt.Errorf("%s: invalid amount of fields in rule %d, got %d want %d", parseErr, i+1, s, fNum)
		}

		p := strings.Split(pohm[pIdx], valuesDlm) // paths

		// TODO: move paths, origin, headers and methods getters to separate methods
		// origins
		var o []string
		if pohm[oIdx] != "" {
			o = strings.Split(pohm[oIdx], valuesDlm)
		}

		// headers
		var h []string
		if pohm[hIdx] != "" {
			h = strings.Split(pohm[hIdx], valuesDlm)
		}

		var m []string // methods
		if pohm[mIdx] == wildcard {
			m = allMethods
		} else if pohm[mIdx] != "" {
			m = strings.Split(strings.ToUpper(pohm[mIdx]), valuesDlm)
			for _, a := range m {
				if ok := contains(validMethods, a); !ok {
					return fmt.Errorf("%s: invalid HTTP method %s in rule %d", parseErr, a, i+1)
				}
			}
		}

		// TODO: check whether path already registered
		for _, v := range p {
			if v == "" {
				return fmt.Errorf("%s: path cannot be empty", parseErr)
			}

			r.pr[v] = rule{
				o: o,
				h: h,
				m: m,
			}

			if v == wildcard {
				// stop parsing when found path wildcard
				return nil
			}
		}
	}

	return nil
}

func contains(l []string, x string) bool {
	for _, a := range l {
		if a == x {
			return true
		}
	}
	return false
}
