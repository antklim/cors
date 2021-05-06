package cors

import (
	"errors"
	"net/http"
	"strings"
)

type rules struct {
	raw string
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

func Routes(paths []string, config string) (http.Handler, error) {
	if len(paths) == 0 {
		return nil, errors.New("invalid paths list: cannot be empty")
	}

	r := newRules(config)

	if err := r.Validate(); err != nil {
		return nil, err
	}

	return nil, nil
}
