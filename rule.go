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
			m = strings.Split(pohm[mIdx], valuesDlm)
			m = toUpper(m)
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

func toUpper(slist []string) []string {
	var a []string
	for _, s := range slist {
		a = append(a, strings.ToUpper(s))
	}
	return a
}
