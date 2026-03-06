package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestHealthEndpoint(t *testing.T) {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	if rr.Body.String() != "ok" {
		t.Fatalf("unexpected body: %s", rr.Body.String())
	}
}