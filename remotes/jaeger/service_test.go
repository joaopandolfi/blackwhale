package jaeger_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joaopandolfi/blackwhale/v2/remotes/jaeger"
)

func TestCreateSpan(t *testing.T) {
	sink := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer sink.Close()
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", sink.URL)

	_, closer := jaeger.Init("teste")
	defer closer.Close()

	handler := func(w http.ResponseWriter, r *http.Request) {
		_, span := jaeger.StartSpanFromRequest(r, "test")
		defer span.End()
		_, child := jaeger.SpanTrace(r.Context(), "child", map[string]any{"k": "v"})
		defer child.End()

		w.Write([]byte(fmt.Sprintf("%s -> %s", "test", "X")))
	}

	rec := httptest.NewRecorder()
	handler(rec, httptest.NewRequest(http.MethodGet, "/health", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}
