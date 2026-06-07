package server

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHealthzProbe(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	srv := Setup("8080", logger, nil) // passing nil for mgr is safe in our modified Setup for tests if we guard it

	ts := httptest.NewServer(srv.Handler)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/healthz")
	if err != nil {
		t.Fatalf("Failed to GET /healthz: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "ok" {
		t.Errorf("Expected body 'ok', got '%v'", string(body))
	}
}
