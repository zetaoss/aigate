package gateway

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"aigate/internal/config"
)

func TestLoggingMiddleware_SkipsHealthz(t *testing.T) {
	g := &Gateway{
		cfg: &config.Config{
			Server: config.ServerConfig{
				LogLevel: "info",
			},
		},
	}

	var buf bytes.Buffer
	originalOutput := log.Writer()
	originalFlags := log.Flags()
	log.SetOutput(&buf)
	log.SetFlags(0)
	t.Cleanup(func() {
		log.SetOutput(originalOutput)
		log.SetFlags(originalFlags)
	})

	handler := LoggingMiddleware(g, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	handler.ServeHTTP(httptest.NewRecorder(), req)

	if strings.Contains(buf.String(), "/healthz") {
		t.Fatalf("expected /healthz request to be excluded from access logs, got: %q", buf.String())
	}
}

func TestLoggingMiddleware_LogsNonHealthz(t *testing.T) {
	g := &Gateway{
		cfg: &config.Config{
			Server: config.ServerConfig{
				LogLevel: "info",
			},
		},
	}

	var buf bytes.Buffer
	originalOutput := log.Writer()
	originalFlags := log.Flags()
	log.SetOutput(&buf)
	log.SetFlags(0)
	t.Cleanup(func() {
		log.SetOutput(originalOutput)
		log.SetFlags(originalFlags)
	})

	handler := LoggingMiddleware(g, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/v1/models", nil)
	handler.ServeHTTP(httptest.NewRecorder(), req)

	if !strings.Contains(buf.String(), "GET /v1/models") {
		t.Fatalf("expected non-healthz request to be access-logged, got: %q", buf.String())
	}
}
