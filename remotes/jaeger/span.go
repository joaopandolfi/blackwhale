package jaeger

import (
	"context"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// StartSpanFromRequest extracts the parent span context from the inbound HTTP request
// and starts a new child span if there is a parent span.
func StartSpanFromRequest(tracer opentracing.Tracer, r *http.Request, name string) (context.Context, opentracing.Span) {
	spanCtx, _ := Extract(tracer, r)
	newSpan := tracer.StartSpan(name, ext.RPCServerOption(spanCtx))
	newCtx := opentracing.ContextWithSpan(r.Context(), newSpan)
	return newCtx, newSpan
}

// SpanTrace creates a tracing and returns the new context and finisher
func SpanTrace(ctx context.Context, operationName string, tags map[string]interface{}) (context.Context, opentracing.Span) {
	// Get span parent
	var parent opentracing.SpanContext
	currentSpan := opentracing.SpanFromContext(ctx)
	if currentSpan != nil {
		parent = currentSpan.Context()
	}
	parentReference := opentracing.ChildOf(parent)

	// Create new span
	newSpan := opentracing.StartSpan(operationName, parentReference, opentracing.Tags(tags))
	// Get context of new span
	newCtx := opentracing.ContextWithSpan(ctx, newSpan)

	return newCtx, newSpan
}
