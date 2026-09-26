package jaeger

import (
	"context"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelprop "go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const tracerName = "blackwhale"

// StartSpanFromRequest extracts the W3C trace context from the inbound HTTP
// request and starts a server span as its child.
func StartSpanFromRequest(r *http.Request, name string) (context.Context, trace.Span) {
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), otelprop.HeaderCarrier(r.Header))
	return otel.Tracer(tracerName).Start(ctx, name, trace.WithSpanKind(trace.SpanKindServer))
}

// SpanTrace starts a client span as a child of the span carried in ctx
func SpanTrace(ctx context.Context, operationName string, tags map[string]any) (context.Context, trace.Span) {
	return otel.Tracer(tracerName).Start(ctx, operationName, trace.WithAttributes(toAttributes(tags)...))
}

func toAttributes(tags map[string]any) []attribute.KeyValue {
	attrs := make([]attribute.KeyValue, 0, len(tags))
	for k, v := range tags {
		switch value := v.(type) {
		case string:
			attrs = append(attrs, attribute.String(k, value))
		case bool:
			attrs = append(attrs, attribute.Bool(k, value))
		case int:
			attrs = append(attrs, attribute.Int(k, value))
		case int64:
			attrs = append(attrs, attribute.Int64(k, value))
		case float64:
			attrs = append(attrs, attribute.Float64(k, value))
		default:
			attrs = append(attrs, attribute.String(k, fmt.Sprint(value)))
		}
	}
	return attrs
}
