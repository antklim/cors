package cors_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/antklim/cors"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOptionsRoutesValidation(t *testing.T) {
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
			_, err := cors.OptionsRoutes(tC.paths, tC.rules)
			assert.EqualError(t, err, tC.err)
		})
	}
}

func TestOptionsRoutes(t *testing.T) {
	paths := []string{"/a", "/b"}

	testCases := []struct {
		desc          string
		method        string
		path          string
		rules         string
		headers       map[string]string
		code          int
		assertHeaders func(*testing.T, http.Header)
	}{
		{
			desc:          "request to unregistered path returns 404",
			method:        http.MethodGet,
			path:          "/not-found",
			rules:         "*;*;;*",
			code:          http.StatusNotFound,
			assertHeaders: func(t *testing.T, h http.Header) {},
		},
		{
			desc:          "non OPTIONS request to registered path returns 405",
			method:        http.MethodGet,
			path:          "/a",
			rules:         "*;*;;*",
			code:          http.StatusMethodNotAllowed,
			assertHeaders: func(t *testing.T, h http.Header) {},
		},
		{
			desc:   "request to registered path A returns 200 for *;*;;* rule",
			method: http.MethodOptions,
			headers: map[string]string{
				"Origin":                        "https://foo.bar.org",
				"Access-Control-Request-Method": "DELETE",
			},
			path:  "/a",
			rules: "*;*;;*",
			code:  http.StatusOK,
			assertHeaders: func(t *testing.T, h http.Header) {
				assert.Equal(t, "*", h.Get("Access-Control-Allow-Origin"))
				assert.Empty(t, h.Get("Access-Control-Allow-Headers"))
				assert.Equal(t, "DELETE", h.Get("Access-Control-Allow-Methods"))
			},
		},
		{
			desc:   "request to registered path B returns 200 for *;*;;* rule",
			method: http.MethodOptions,
			headers: map[string]string{
				"Origin":                        "https://foo.bar.org",
				"Access-Control-Request-Method": "PUT",
			},
			path:  "/b",
			rules: "*;*;;*",
			code:  http.StatusOK,
			assertHeaders: func(t *testing.T, h http.Header) {
				assert.Equal(t, "*", h.Get("Access-Control-Allow-Origin"))
				assert.Empty(t, h.Get("Access-Control-Allow-Headers"))
				assert.Equal(t, "PUT", h.Get("Access-Control-Allow-Methods"))
			},
		},
		{
			desc:   "request to unregistered header returns 403",
			method: http.MethodOptions,
			headers: map[string]string{
				"Origin":                         "https://foo.bar.org",
				"Access-Control-Request-Headers": "content-type",
				"Access-Control-Request-Method":  "PUT",
			},
			path:          "/a",
			rules:         "*;*;;*",
			code:          http.StatusForbidden,
			assertHeaders: func(t *testing.T, h http.Header) {},
		},
		{
			desc:   "request to custom configured path /a returns 200",
			method: http.MethodOptions,
			headers: map[string]string{
				"Origin":                         "https://foo.bar.org",
				"Access-Control-Request-Headers": "content-type",
				"Access-Control-Request-Method":  "PUT",
			},
			path:  "/a",
			rules: "/a;https://foo.bar.org;content-type,x-correlation-id;PUT\n*;bar.com;;*",
			code:  http.StatusOK,
			assertHeaders: func(t *testing.T, h http.Header) {
				assert.Equal(t, "https://foo.bar.org", h.Get("Access-Control-Allow-Origin"))
				assert.Equal(t, "Content-Type", h.Get("Access-Control-Allow-Headers"))
				assert.Equal(t, "PUT", h.Get("Access-Control-Allow-Methods"))
			},
		},
		{
			desc:   "request to custom configured path /b returns 200",
			method: http.MethodOptions,
			headers: map[string]string{
				"Origin":                        "https://bar.foo.org",
				"Access-Control-Request-Method": "DELETE",
			},
			path:  "/b",
			rules: "/a;https://foo.bar.org;content-type,x-correlation-id;PUT\n*;;;*",
			code:  http.StatusOK,
			assertHeaders: func(t *testing.T, h http.Header) {
				assert.Equal(t, "*", h.Get("Access-Control-Allow-Origin"))
				assert.Empty(t, h.Get("Access-Control-Allow-Headers"))
				assert.Equal(t, "DELETE", h.Get("Access-Control-Allow-Methods"))
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			req := httptest.NewRequest(tC.method, tC.path, nil)

			for k, v := range tC.headers {
				req.Header.Set(k, v)
			}

			rr := httptest.NewRecorder()

			h, err := cors.OptionsRoutes(paths, tC.rules)
			require.NoError(t, err)
			h.ServeHTTP(rr, req)

			res := rr.Result()
			assert.Equal(t, tC.code, res.StatusCode)
			tC.assertHeaders(t, res.Header)
		})
	}
}

func TestRouteMiddleware(t *testing.T) {
	path := "/a"
	mainHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	})

	config := "*;https://foo.bar.org;content-type;*"
	rules := cors.NewRules(config)
	err := rules.Parse()
	require.NoError(t, err)

	rule, ok := rules.OfPath(path)
	require.True(t, ok)

	h := cors.Middleware("*", rule)(mainHandler)

	router := mux.NewRouter()
	router.Handle(path, h).Methods(http.MethodOptions, http.MethodGet, http.MethodPut)

	// OPTIONS for GET request
	{
		req := httptest.NewRequest(http.MethodOptions, path, nil)
		rr := httptest.NewRecorder()

		req.Header.Set("Origin", "https://foo.bar.org")
		req.Header.Set("Access-Control-Request-Method", "GET")

		router.ServeHTTP(rr, req)
		res := rr.Result()

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "https://foo.bar.org", res.Header.Get("Access-Control-Allow-Origin"))
		assert.Empty(t, res.Header.Get("Access-Control-Allow-Methods")) // GET is allowed by default
	}
	// GET request
	{
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rr := httptest.NewRecorder()

		req.Header.Set("Origin", "https://foo.bar.org")

		router.ServeHTTP(rr, req)
		res := rr.Result()
		defer res.Body.Close()

		body, _ := ioutil.ReadAll(res.Body)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "https://foo.bar.org", res.Header.Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "OK", string(body))
	}

	// OPTIONS for PUT request
	{
		req := httptest.NewRequest(http.MethodOptions, path, nil)
		rr := httptest.NewRecorder()

		req.Header.Set("Origin", "https://foo.bar.org")
		req.Header.Set("Access-Control-Request-Headers", "content-type")
		req.Header.Set("Access-Control-Request-Method", "PUT")

		router.ServeHTTP(rr, req)
		res := rr.Result()

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "https://foo.bar.org", res.Header.Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "Content-Type", res.Header.Get("Access-Control-Allow-Headers"))
		assert.Equal(t, "PUT", res.Header.Get("Access-Control-Allow-Methods"))
	}
	// PUT request
	{
		req := httptest.NewRequest(http.MethodPut, path, nil)
		rr := httptest.NewRecorder()

		req.Header.Set("Origin", "https://foo.bar.org")
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, req)
		res := rr.Result()
		defer res.Body.Close()

		body, _ := ioutil.ReadAll(res.Body)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "https://foo.bar.org", res.Header.Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "OK", string(body))
	}

	// OPTIONS request without method
	{
		req := httptest.NewRequest(http.MethodOptions, path, nil)
		rr := httptest.NewRecorder()

		req.Header.Set("Origin", "https://foo.bar.org")
		req.Header.Set("Access-Control-Request-Headers", "X-Correlation-Id")

		router.ServeHTTP(rr, req)
		res := rr.Result()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	}

	// OPTIONS request witho unregistered header
	{
		req := httptest.NewRequest(http.MethodOptions, path, nil)
		rr := httptest.NewRecorder()

		req.Header.Set("Origin", "https://foo.bar.org")
		req.Header.Set("Access-Control-Request-Headers", "X-Correlation-Id")
		req.Header.Set("Access-Control-Request-Method", "PUT")

		router.ServeHTTP(rr, req)
		res := rr.Result()

		assert.Equal(t, http.StatusForbidden, res.StatusCode)
	}
}
