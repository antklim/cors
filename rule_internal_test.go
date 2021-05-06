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
		// {
		// TODO: add fields amount validation
		// 	desc: "fails when cors rules config has invalid amount fields in a rule",
		// },
		// {
		// TODO: add method validation
		// 	desc: "fails when cors rules config has invalid http method",
		// },
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
		{
			desc:   "parses config with explicit fields",
			config: "/a,/b;foo.com,bar.com;content-type,content-length;DELETE,PUT",
			r: &rules{
				raw: "/a,/b;foo.com,bar.com;content-type,content-length;DELETE,PUT",
				r: []rule{{
					raw: "/a,/b;foo.com,bar.com;content-type,content-length;DELETE,PUT",
					p:   []string{"/a", "/b"},
					o:   []string{"foo.com", "bar.com"},
					h:   []string{"content-type", "content-length"},
					m:   []string{http.MethodDelete, http.MethodPut},
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
