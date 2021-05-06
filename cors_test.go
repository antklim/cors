package cors_test

import (
	"testing"

	"github.com/antklim/cors"
	"github.com/stretchr/testify/assert"
)

// TODO: Routes(<list of paths>, <cors config>) http.Handler
//			 cors config: paths;allowed origins;allowed headers;allowed methods\n
// path can be *
// allowed origins can be *
// allowed headers can be explicit

func TestRoutesValidation(t *testing.T) {
	testCases := []struct {
		desc string
		p    []string
		r    string
		err  string
	}{
		{
			desc: "fails when paths list is empty",
			err:  "invalid paths list: cannot be empty",
		},
		{
			desc: "fails when cors rules is empty",
			p:    []string{"/a"},
			err:  "invalid cors rules: cannot be empty",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			_, err := cors.Routes(tC.p, tC.r)
			assert.EqualError(t, err, tC.err)
		})
	}
}
