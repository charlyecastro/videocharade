package handlers

import (
	"final-project-crew/videoCharade/servers/gateway/models/users"
	"final-project-crew/videoCharade/servers/gateway/sessions"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func setUpServer() *httptest.Server {
	ctx := NewContext("test",
		sessions.NewMemStore(time.Hour, time.Minute),
		users.NewMockStore())
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/users", ctx.UsersHandler)

	ts := httptest.NewServer(NewCors(mux))
	return ts
}

func TestBadOptionsRequest(t *testing.T) {
	const expectedAllowOrigin = "*"
	const expectedAllowMethods = "GET, PUT, POST, PATCH, DELETE"
	const expectedAllowHeaders = "Content-Type, Authorization"
	const expectedExposed = "Authorization"
	const expectedMaxAge = "Access-Control-Max-Age: 600"

	ts := setUpServer()
	defer ts.Close()
	client := ts.Client()
	req, _ := http.NewRequest("OPTIONS", ts.URL, nil)
	res, _ := client.Do(req)
	header := res.Header

	allowOrigin := header.Get("Access-Control-Allow-Origin")
	if allowOrigin != expectedAllowOrigin {
		t.Errorf("incorrect status code: expected %s but got %s",
			allowOrigin, expectedAllowOrigin)
	}

	allowMethods := header.Get("Access-Control-Allow-Methods")
	if allowMethods != expectedAllowMethods {
		t.Errorf("incorrect status code: expected %s but got %s",
			allowMethods, expectedAllowMethods)
	}

	allowHeaders := header.Get("Access-Control-Allow-Headers")
	if allowHeaders != expectedAllowHeaders {
		t.Errorf("incorrect status code: expected %s but got %s",
			allowHeaders, expectedAllowHeaders)
	}

	exposed := header.Get("Access-Control-Expose-Headers")
	if exposed != expectedExposed {
		t.Errorf("incorrect status code: expected %s but got %s",
			exposed, expectedExposed)
	}

	maxAge := header.Get("Access-Control-Max-Age")
	if maxAge != expectedMaxAge {
		t.Errorf("incorrect status code: expected %s but got %s",
			maxAge, expectedMaxAge)
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("incorrect status code: expected %d but got %d",
			res.StatusCode, http.StatusBadRequest)
	}

}

func TestGoodOptionsRequest(t *testing.T) {
	const requestMethod = "Access-Control-Request-Method"
	const origin = "Origin"

	ts := setUpServer()
	defer ts.Close()
	client := ts.Client()
	req, _ := http.NewRequest("OPTIONS", ts.URL, nil)
	req.Header.Set(requestMethod, http.MethodGet)
	req.Header.Set(origin, "localhost")
	res, _ := client.Do(req)

	if res.StatusCode != http.StatusOK {
		t.Errorf("incorrect status code: expected %d but got %d",
			res.StatusCode, http.StatusBadRequest)
	}

}
