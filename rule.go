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

type exprType string

const (
	ruleOrigins exprType = "origins"
	ruleHeaders exprType = "headers"
	ruleMethods exprType = "methods"
)

type Rule struct {
	o []string // origins
	h []string // headers
	m []string // methods
}

func (r Rule) Origins() []string {
	return r.o
}

func (r Rule) Headers() []string {
	return r.h
}

func (r Rule) Methods() []string {
	return r.m
}

type RuleBuilder struct {
	expr map[exprType][]string
}

func NewRuleBuilder() RuleBuilder {
	return RuleBuilder{}
}

func (b RuleBuilder) WithOrigins(o ...string) RuleBuilder {
	if len(o) > 0 {
		if b.expr == nil {
			b.expr = make(map[exprType][]string)
		}
		b.expr[ruleOrigins] = o
	}
	return b
}

func (b RuleBuilder) WithHeaders(h ...string) RuleBuilder {
	if len(h) > 0 {
		if b.expr == nil {
			b.expr = make(map[exprType][]string)
		}
		b.expr[ruleHeaders] = h
	}
	return b
}

func (b RuleBuilder) WithMethods(m ...string) RuleBuilder {
	vm := filterMethods(m)
	if len(vm) > 0 {
		if b.expr == nil {
			b.expr = make(map[exprType][]string)
		}
		b.expr[ruleMethods] = vm
	}
	return b
}

func (b RuleBuilder) Build() Rule {
	r := Rule{}
	for k, v := range b.expr {
		switch k {
		case ruleOrigins:
			r.o = v
		case ruleHeaders:
			r.h = v
		case ruleMethods:
			r.m = v
		}
	}
	return r
}

type Rules struct {
	raw string
	op  []string // ordered paths list
	pr  map[string]Rule
}

func NewRules(config string) *Rules {
	return &Rules{raw: config}
}

func (r *Rules) Parse() error {
	return r.parseTxt()
}

func (r *Rules) Paths() []string {
	return r.op
}

func (r *Rules) OfPath(path string) (Rule, bool) {
	if rule, ok := r.pr[path]; ok {
		return rule, true
	}

	if rule, ok := r.pr[wildcard]; ok {
		return rule, true
	}

	return Rule{}, false
}

func (r *Rules) parseTxt() error {
	if strings.TrimSpace(r.raw) == "" {
		return fmt.Errorf("%s: cannot be empty", parseErr)
	}

	rawRules := strings.Split(r.raw, rulesDlm)

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

			if r.op == nil {
				r.op = append([]string{}, p)
			} else {
				r.op = append(r.op, p)
			}

			if r.pr == nil {
				r.pr = make(map[string]Rule)
			}

			r.pr[p] = Rule{
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

func filterMethods(mm []string) []string {
	fm := make([]string, 0, len(mm))
	for _, m := range mm {
		if contains(validMethods, m) {
			fm = append(fm, m)
		}
	}
	return fm
}
