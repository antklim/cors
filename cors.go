package cors

import (
	"errors"
	"net/http"
	"strings"
)

func Routes(paths []string, rules string) (http.Handler, error) {
	if len(paths) == 0 {
		return nil, errors.New("invalid paths list: cannot be empty")
	}

	if strings.TrimSpace(rules) == "" {
		return nil, errors.New("invalid cors rules: cannot be empty")
	}

	return nil, nil
}
