package cors

import (
	"errors"
	"strings"
)

// const (
// 	rulesDlm  = '\n'
// 	fieldsDlm = ';'
// )

// type rule struct {
// 	raw string
// 	p   []string // paths
// 	o   []string // origins
// 	h   []string // headers
// 	m   []string // methods
// }

type rules struct {
	raw string
	// r   []rule
}

func newRules(config string) *rules {
	return &rules{raw: config}
}

func (r *rules) Parse() error {
	if strings.TrimSpace(r.raw) == "" {
		return errors.New("invalid cors rules: cannot be empty")
	}
	return nil
}
