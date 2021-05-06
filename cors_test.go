package cors_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/antklim/cors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO: Routes(<list of paths>, <cors config>) http.Handler
//			 cors config: paths;allowed origins;allowed headers;allowed methods\n
// path can be *
// allowed origins can be *
// allowed headers can be explicit

func TestRoutesValidation(t *testing.T) {
	testCases := []struct {
		desc  string
		paths []string
		rules string
		err   string
	}{
		{
			desc: "fails when paths list is empty",
			err:  "invalid paths list: cannot be empty",
		},
		{
			desc:  "fails when cors rules is empty",
			paths: []string{"/a"},
			err:   "invalid cors rules: cannot be empty",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			_, err := cors.Routes(tC.paths, tC.rules)
			assert.EqualError(t, err, tC.err)
		})
	}
}

// [/a, /b] *;*;;* - OPTIONS request to a and b returns allowed origin * and all allowed methods
// [/a, /b] /a;example.com;content-type;DELETE - OPTIONS request to a returns allowed origin example.com and allowed method
// [/a, /b] /a;example.com;content-type;DELETE - OPTIONS request to b returns 405
// [/a, /b] /a;example.com;content-type;DELETE\n/b;example.com;content-type;POST,PUT - OPTIONS request to b returns 200

func TestRoutes(t *testing.T) {
	paths := []string{"/a", "/b"}

	testCases := []struct {
		desc   string
		method string
		path   string
		rules  string
		code   int
	}{
		{
			desc:   "request to unregistered path returns 404",
			method: http.MethodGet,
			path:   "/not-found",
			rules:  "*;*;;*",
			code:   http.StatusNotFound,
		},
		{
			desc:   "non OPTIONS request to registered path returns 405",
			method: http.MethodGet,
			path:   "/a",
			rules:  "*;*;;*",
			code:   http.StatusMethodNotAllowed,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			req := httptest.NewRequest(tC.method, tC.path, nil)
			rr := httptest.NewRecorder()

			h, err := cors.Routes(paths, tC.rules)
			require.NoError(t, err)
			h.ServeHTTP(rr, req)

			res := rr.Result()
			assert.Equal(t, tC.code, res.StatusCode)
		})
	}
}
