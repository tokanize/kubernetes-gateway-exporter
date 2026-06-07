package server

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// A mock cache that we can control for the readyz endpoint
type mockCache struct {
	synced bool
}

func (m *mockCache) WaitForCacheSync(ctx context.Context) bool {
	return m.synced
}

// A mock manager wrapping the mock cache
type mockManager struct {
	cache *mockCache
}

func (m *mockManager) GetCache() interface {
	WaitForCacheSync(ctx context.Context) bool
} {
	return m.cache
}

// Add any other required methods for manager.Manager here returning nil/panicking if needed,
// but because Go interfaces are satisfied implicitly, we only need what's used if we typecast.
// Actually, our Setup takes manager.Manager. Since we only call GetCache().WaitForCacheSync,
// we'll need to implement the entire manager.Manager interface or just pass nil.
// To keep things simple and table-driven without huge mocks, we'll test the endpoints via the Server directly.

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
