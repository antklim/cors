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
	fields    int    = 4

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

func newRules(config string) *rules {
	return &rules{raw: config}
}

func (r *rules) Parse() error {
	if strings.TrimSpace(r.raw) == "" {
		return errors.New("invalid cors rules: cannot be empty")
	}

	rawRules := strings.Split(r.raw, rulesDlm)
	r.r = make([]rule, len(rawRules))
	for i, rr := range rawRules {
		pohm := strings.Split(rr, fieldsDlm)

		if s := len(pohm); s != fields {
			return fmt.Errorf("invalid cors rules: invalid amount of fields in rule %d, got %d want %d", i+1, s, fields)
		}

		p := []string{pohm[pIdx]} // paths
		o := []string{pohm[oIdx]} // origins
		h := []string{pohm[hIdx]} // headers
		m := []string{pohm[mIdx]} // methods

		if pohm[mIdx] == wildcard {
			m = allMethods
		}

		r.r[i] = rule{
			raw: rr,
			p:   p,
			o:   o,
			h:   h,
			m:   m,
		}
	}

	return nil
}
