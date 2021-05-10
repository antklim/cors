package cors

import (
	"fmt"
	"net/http"
	"strings"
)

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
	return r.parseTxt()
}

func (r *rules) parseTxt() error {
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

		paths := parsePaths(pohm[pIdx])
		if paths == nil {
			return fmt.Errorf("%s: path cannot be empty", parseErr)
		}

		origins := parseOrigins(pohm[oIdx])
		headers := parseHeaders(pohm[hIdx])
		methods, err := parseMethods(pohm[mIdx], i+1)
		if err != nil {
			return err
		}

		for _, p := range paths {
			if p == "" {
				return fmt.Errorf("%s: path cannot be empty", parseErr)
			}

			// ignore repeatable occurrences of path in config
			if _, ok := r.pr[p]; ok {
				continue
			}

			r.pr[p] = rule{
				o: origins,
				h: headers,
				m: methods,
			}

			if p == wildcard {
				// stop parsing when found path wildcard
				return nil
			}
		}
	}

	return nil
}

func parsePaths(s string) []string {
	var p []string
	if s != "" {
		p = strings.Split(s, valuesDlm)
	}
	return p
}

func parseOrigins(s string) []string {
	var o []string
	if s != "" {
		o = strings.Split(s, valuesDlm)
	}
	return o
}

func parseHeaders(s string) []string {
	var h []string
	if s != "" {
		h = strings.Split(s, valuesDlm)
	}
	return h
}

func parseMethods(s string, ruleNum int) ([]string, error) {
	s = strings.TrimSpace(s)

	if s == "" {
		return nil, nil
	}

	if s == wildcard {
		return allMethods, nil
	}

	m := strings.Split(strings.ToUpper(s), valuesDlm)
	for _, a := range m {
		if ok := contains(validMethods, a); !ok {
			return nil, fmt.Errorf("%s: invalid HTTP method %s in rule %d", parseErr, a, ruleNum)
		}
	}

	return m, nil
}

func contains(l []string, x string) bool {
	for _, a := range l {
		if a == x {
			return true
		}
	}
	return false
}
