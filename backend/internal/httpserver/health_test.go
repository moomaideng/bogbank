package httpserver

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	srv := New(":0", "test", nil)

	t.Run("livez", func(t *testing.T) {
		rec := request(t, srv.Handler(), http.MethodGet, "/livez")
		if rec.Code != http.StatusOK {
			t.Fatalf("status %d", rec.Code)
		}
	})

	t.Run("readyz without dependencies", func(t *testing.T) {
		rec := request(t, srv.Handler(), http.MethodGet, "/readyz")
		if rec.Code != http.StatusOK {
			t.Fatalf("status %d", rec.Code)
		}
	})

	t.Run("openapi", func(t *testing.T) {
		rec := request(t, srv.Handler(), http.MethodGet, "/openapi.json")
		if rec.Code != http.StatusOK {
			t.Fatalf("status %d body %s", rec.Code, rec.Body.String())
		}
	})
}

func TestReadyzFailsWhenPingerFails(t *testing.T) {
	srv := New(":0", "test", []Pinger{failPing{}})
	rec := request(t, srv.Handler(), http.MethodGet, "/readyz")
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status %d", rec.Code)
	}
}

type failPing struct{}

func (failPing) PingContext(context.Context) error {
	return errors.New("down")
}

func request(t *testing.T, handler http.Handler, method, path string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, nil)
	handler.ServeHTTP(rec, req)
	return rec
}
