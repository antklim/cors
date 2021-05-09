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

		if len(rr) == 0 {
			continue
		}

		pohm := strings.Split(rr, fieldsDlm)

		if s := len(pohm); s != fNum {
			return fmt.Errorf("%s: invalid amount of fields in rule %d, got %d want %d", parseErr, i+1, s, fNum)
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
		for _, v := range p {
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
