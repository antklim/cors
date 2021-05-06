package cors

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRuleParseError(t *testing.T) {
	testCases := []struct {
		desc   string
		config string
		err    string
	}{
		{
			desc: "fails when cors rules config is empty",
			err:  "invalid cors rules: cannot be empty",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			rules := newRules(tC.config)
			err := rules.Parse()
			assert.EqualError(t, err, tC.err)
		})
	}
}

func TestRuleParse(t *testing.T) {
	testCases := []struct {
		desc   string
		config string
		r      *rules
	}{
		{
			desc:   "parses wildcard config",
			config: "*;*;;*",
			r: &rules{
				raw: "*;*;;*",
				r: []rule{{
					raw: "*;*;;*",
					p:   []string{"*"},
					o:   []string{"*"},
					h:   []string{""},
					m: []string{
						http.MethodDelete, http.MethodGet, http.MethodHead,
						http.MethodPatch, http.MethodPost, http.MethodPut,
					},
				}},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			rules := newRules(tC.config)
			err := rules.Parse()
			require.NoError(t, err)
			assert.Equal(t, tC.r, rules)
		})
	}
}
