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
		{
			desc:   "fails when path is empty",
			config: ";;;",
			err:    "invalid cors rules: path cannot be empty",
		},
		{
			desc:   "fails when cors rules config has invalid amount fields in a rule",
			config: "*;;",
			err:    "invalid cors rules: invalid amount of fields in rule 1, got 3 want 4",
		},
		{
			desc:   "fails when cors rules config has invalid http method",
			config: "*;;;foo",
			err:    "invalid cors rules: invalid HTTP method FOO in rule 1",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			rules := NewRules(tC.config)
			err := rules.Parse()
			assert.EqualError(t, err, tC.err)
		})
	}
}

func TestRuleParse(t *testing.T) {
	testCases := []struct {
		desc   string
		config string
		r      *Rules
	}{
		{
			desc:   "parses wildcard config",
			config: "*;*;;*",
			r: &Rules{
				raw: "*;*;;*",
				pr: map[string]Rule{
					"*": {
						o: []string{"*"},
						h: nil,
						m: []string{
							http.MethodDelete, http.MethodGet, http.MethodHead,
							http.MethodPatch, http.MethodPost, http.MethodPut,
						},
					},
				},
			},
		},
		{
			desc:   "parses empty origin, headers and methods",
			config: "*;;;",
			r: &Rules{
				raw: "*;;;",
				pr: map[string]Rule{
					"*": {
						o: nil,
						h: nil,
						m: nil,
					},
				},
			},
		},
		{
			desc:   "parses any case methods",
			config: "*;;;post,Put",
			r: &Rules{
				raw: "*;;;post,Put",
				pr: map[string]Rule{
					"*": {
						o: nil,
						h: nil,
						m: []string{http.MethodPost, http.MethodPut},
					},
				},
			},
		},
		{
			desc:   "parses config with explicit fields",
			config: "/a,/b;foo.com,bar.com;content-type,content-length;DELETE,PUT",
			r: &Rules{
				raw: "/a,/b;foo.com,bar.com;content-type,content-length;DELETE,PUT",
				pr: map[string]Rule{
					"/a": {
						o: []string{"foo.com", "bar.com"},
						h: []string{"content-type", "content-length"},
						m: []string{http.MethodDelete, http.MethodPut},
					},
					"/b": {
						o: []string{"foo.com", "bar.com"},
						h: []string{"content-type", "content-length"},
						m: []string{http.MethodDelete, http.MethodPut},
					},
				},
			},
		},
		{
			desc: "parses multiline config",
			config: `/a;foo.com;content-type;DELETE
			/b;bar.com;content-length;PUT`,
			r: &Rules{
				raw: `/a;foo.com;content-type;DELETE
			/b;bar.com;content-length;PUT`,
				pr: map[string]Rule{
					"/a": {
						o: []string{"foo.com"},
						h: []string{"content-type"},
						m: []string{http.MethodDelete},
					},
					"/b": {
						o: []string{"bar.com"},
						h: []string{"content-length"},
						m: []string{http.MethodPut},
					},
				},
			},
		},
		{
			desc: "ignores leading and closing empty strings in multiline config",
			config: `
					/a;foo.com;content-type;DELETE
					/b;bar.com;content-length;PUT
					`,
			r: &Rules{
				raw: `
					/a;foo.com;content-type;DELETE
					/b;bar.com;content-length;PUT
					`,
				pr: map[string]Rule{
					"/a": {
						o: []string{"foo.com"},
						h: []string{"content-type"},
						m: []string{http.MethodDelete},
					},
					"/b": {
						o: []string{"bar.com"},
						h: []string{"content-length"},
						m: []string{http.MethodPut},
					},
				},
			},
		},
		{
			desc: "ignores repeatable occurrences of path in config",
			config: `/a;foo.com;content-type;DELETE
			/a;bar.com;content-length;PUT`,
			r: &Rules{
				raw: `/a;foo.com;content-type;DELETE
			/a;bar.com;content-length;PUT`,
				pr: map[string]Rule{
					"/a": {
						o: []string{"foo.com"},
						h: []string{"content-type"},
						m: []string{http.MethodDelete},
					},
				},
			},
		},
		{
			desc: "stops parsing when found paths wildcard",
			config: `/a;foo.com;content-type;DELETE
			*;foobar.com;;PATCH
			/b;bar.com;content-length;PUT`,
			r: &Rules{
				raw: `/a;foo.com;content-type;DELETE
			*;foobar.com;;PATCH
			/b;bar.com;content-length;PUT`,
				pr: map[string]Rule{
					"/a": {
						o: []string{"foo.com"},
						h: []string{"content-type"},
						m: []string{http.MethodDelete},
					},
					"*": {
						o: []string{"foobar.com"},
						h: nil,
						m: []string{http.MethodPatch},
					},
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			rules := NewRules(tC.config)
			err := rules.Parse()
			require.NoError(t, err)
			assert.Equal(t, tC.r, rules)
		})
	}
}
