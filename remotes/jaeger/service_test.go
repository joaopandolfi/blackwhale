package jaeger_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/joaopandolfi/blackwhale/remotes/jaeger"
	"github.com/opentracing/opentracing-go"
)

var handler func(w http.ResponseWriter, r *http.Request)

func TestCreateSpan(t *testing.T) {
	tracer, closer := jaeger.Init("teste")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	handler = func(w http.ResponseWriter, r *http.Request) {
		_, span := jaeger.SpanTrace(r.Context(), "test", map[string]interface{}{})
		defer span.Finish()

		w.Write([]byte(fmt.Sprintf("%s -> %s", "test", "X")))
	}
}
