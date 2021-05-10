package cors_test

import (
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
				// TODO: compare origins, headers, methods
			},
		},
		{
			desc:   "is the wildcard rule otherwise",
			config: "/a;foo.com;content-type;DELETE\n*;foobar.com;;PATCH",
			path:   "/b",
			assert: func(t *testing.T, r cors.Rule, found bool) {
				assert.True(t, found)
				// TODO: compare origins, headers, methods
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

func TestRuleOrigins(t *testing.T) {

}
