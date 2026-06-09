package server

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/tokanize/kubernetes-gateway-exporter/internal/kubernetes/mapper"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func TestHealthzProbe(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	srv := Setup("8080", logger, nil, nil) // passing nil for mgr and mapper is safe in our modified Setup for tests if we guard it

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

func TestReadinessProbe(t *testing.T) {
	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		t.Fatalf("Failed to add corev1 to scheme: %v", err)
	}
	if err := gwv1.Install(scheme); err != nil {
		t.Fatalf("Failed to add gwv1 to scheme: %v", err)
	}
	c := fake.NewClientBuilder().WithScheme(scheme).Build()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	m := mapper.NewMapper(c, logger)

	// Case 1: Mapper is not ready (initially m.Ready() is false)
	srv := Setup("8080", logger, nil, m)
	ts := httptest.NewServer(srv.Handler)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/readyz")
	if err != nil {
		t.Fatalf("Failed to GET /readyz: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("Expected status 503 Service Unavailable, got %v", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "mapper cache not ready" {
		t.Errorf("Expected body 'mapper cache not ready', got '%s'", string(body))
	}

	// Case 2: Mapper is ready (update cache once)
	if err := m.UpdateCache(context.Background()); err != nil {
		t.Fatalf("Failed to update cache: %v", err)
	}

	resp2, err := http.Get(ts.URL + "/readyz")
	if err != nil {
		t.Fatalf("Failed to GET /readyz (ready): %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", resp2.StatusCode)
	}

	body2, _ := io.ReadAll(resp2.Body)
	if string(body2) != "ready" {
		t.Errorf("Expected body 'ready', got '%s'", string(body2))
	}
}
