package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"amazon-vl/internal/auth"
)

func TestNewServer(t *testing.T) {
	cfg := Config{
		Dir:  "/tmp",
		Port: "8080",
		Auth: auth.DefaultConfig(),
	}

	srv := New(cfg)

	if srv == nil {
		t.Error("expected non-nil server")
	}

	if srv.config.Dir != "/tmp" {
		t.Errorf("expected dir '/tmp', got '%s'", srv.config.Dir)
	}

	if srv.config.Port != "8080" {
		t.Errorf("expected port '8080', got '%s'", srv.config.Port)
	}
}

func TestHealthHandler(t *testing.T) {
	cfg := Config{
		Dir:  "/tmp",
		Port: "8080",
		Auth: auth.DefaultConfig(),
	}

	srv := New(cfg)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	srv.healthHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	expected := `{"status":"healthy"}`
	if rec.Body.String() != expected {
		t.Errorf("expected body '%s', got '%s'", expected, rec.Body.String())
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", contentType)
	}
}

func TestNewFileServer(t *testing.T) {
	fs := NewFileServer("/var/log")

	if fs == nil {
		t.Error("expected non-nil file server")
	}

	if fs.Dir != "/var/log" {
		t.Errorf("expected dir '/var/log', got '%s'", fs.Dir)
	}
}
