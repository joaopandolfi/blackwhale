package jaeger_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/joaopandolfi/blackwhale/v2/remotes/jaeger"
	"github.com/opentracing/opentracing-go"
)

var handler func(w http.ResponseWriter, r *http.Request)

func TestCreateSpan(t *testing.T) {
	t.Setenv("JAEGER_SERVICE_NAME", "teste")
	tracer, closer := jaeger.Init("teste")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	handler = func(w http.ResponseWriter, r *http.Request) {
		_, span := jaeger.SpanTrace(r.Context(), "test", map[string]any{})
		defer span.Finish()

		w.Write([]byte(fmt.Sprintf("%s -> %s", "test", "X")))
	}
}
