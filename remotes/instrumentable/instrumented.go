// Package instrumentable provides span-name helpers for tracing.
package instrumentable

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/trace"

	"github.com/joaopandolfi/blackwhale/v2/remotes/jaeger"
)

type Instrumented struct {
	SpanName string
}

func New(name string) Instrumented {
	return Instrumented{
		SpanName: name,
	}
}

func (s *Instrumented) SpanTrace(ctx context.Context, name string, tags map[string]any) (context.Context, trace.Span) {
	return jaeger.SpanTrace(ctx, fmt.Sprintf("%s.%s", s.SpanName, name), tags)
}
