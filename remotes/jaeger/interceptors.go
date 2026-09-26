package jaeger

import (
	"net/http"

	"go.opentelemetry.io/otel"
	otelprop "go.opentelemetry.io/otel/propagation"
)

// Inject propagates the trace context held in request.Context() into the
// outbound request headers (W3C traceparent/tracestate).
func Inject(request *http.Request) {
	otel.GetTextMapPropagator().Inject(request.Context(), otelprop.HeaderCarrier(request.Header))
}
