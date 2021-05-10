package cors_test

import (
	"net/http"
	"testing"

	"github.com/antklim/cors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRulesPaths(t *testing.T) {
	config := `
		/a;foo.com;content-type;DELETE
		*;foobar.com;;PATCH
		/b;bar.com;content-length;PUT`

	rules := cors.NewRules(config)
	err := rules.Parse()
	require.NoError(t, err)
	assert.Equal(t, []string{"/a", "*"}, rules.Paths())
}

func TestRulesOfPath(t *testing.T) {
	testCases := []struct {
		desc   string
		config string
		path   string
		assert func(*testing.T, cors.Rule, bool)
	}{
		{
			desc:   "not found when no path and no wildcard configured",
			config: "/a;foo.com;content-type;DELETE",
			path:   "/b",
			assert: func(t *testing.T, r cors.Rule, found bool) {
				assert.False(t, found)
			},
		},
		{
			desc:   "is the path rule when path explicitly configured",
			config: "/a;foo.com;content-type;DELETE",
			path:   "/a",
			assert: func(t *testing.T, r cors.Rule, found bool) {
				assert.True(t, found)
				assert.Equal(t, []string{"foo.com"}, r.Origins())
				assert.Equal(t, []string{"content-type"}, r.Headers())
				assert.Equal(t, []string{http.MethodDelete}, r.Methods())
			},
		},
		{
			desc:   "is the wildcard rule otherwise",
			config: "/a;foo.com;content-type;DELETE\n*;bar.com;;PATCH",
			path:   "/b",
			assert: func(t *testing.T, r cors.Rule, found bool) {
				assert.True(t, found)
				assert.Equal(t, []string{"bar.com"}, r.Origins())
				assert.Nil(t, r.Headers())
				assert.Equal(t, []string{http.MethodPatch}, r.Methods())
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			rules := cors.NewRules(tC.config)
			err := rules.Parse()
			require.NoError(t, err)
			rule, ok := rules.OfPath(tC.path)
			tC.assert(t, rule, ok)
		})
	}
}

func TestRuleBuilder(t *testing.T) {
	testCases := []struct {
		desc   string
		o      []string
		h      []string
		m      []string
		assert func(*testing.T, cors.Rule)
	}{
		{
			desc: "default build",
			assert: func(t *testing.T, r cors.Rule) {
				assert.Nil(t, r.Origins())
				assert.Nil(t, r.Headers())
				assert.Nil(t, r.Methods())
			},
		},
		{
			desc: "custom origins build",
			o:    []string{"a", "b"},
			assert: func(t *testing.T, r cors.Rule) {
				assert.Equal(t, []string{"a", "b"}, r.Origins())
				assert.Nil(t, r.Headers())
				assert.Nil(t, r.Methods())
			},
		},
		{
			desc: "custom headers build",
			h:    []string{"content-type"},
			assert: func(t *testing.T, r cors.Rule) {
				assert.Nil(t, r.Origins())
				assert.Equal(t, []string{"content-type"}, r.Headers())
				assert.Nil(t, r.Methods())
			},
		},
		{
			desc: "custom methods build",
			m:    []string{http.MethodDelete, http.MethodPatch},
			assert: func(t *testing.T, r cors.Rule) {
				assert.Nil(t, r.Origins())
				assert.Nil(t, r.Headers())
				assert.Equal(t, []string{http.MethodDelete, http.MethodPatch}, r.Methods())
			},
		},
		{
			desc: "fully custom build",
			o:    []string{"a", "b"},
			h:    []string{"content-type"},
			m:    []string{http.MethodDelete, http.MethodPatch},
			assert: func(t *testing.T, r cors.Rule) {
				assert.Equal(t, []string{"a", "b"}, r.Origins())
				assert.Equal(t, []string{"content-type"}, r.Headers())
				assert.Equal(t, []string{http.MethodDelete, http.MethodPatch}, r.Methods())
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			rule := cors.NewRuleBuilder().
				WithOrigins(tC.o...).
				WithHeaders(tC.h...).
				WithMethods(tC.m...).
				Build()
			tC.assert(t, rule)
		})
	}
}
