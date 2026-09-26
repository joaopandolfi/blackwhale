package instrumentable

import (
	"context"
	"fmt"

	"github.com/joaopandolfi/blackwhale/v2/remotes/jaeger"
	"github.com/opentracing/opentracing-go"
)

type Instrumented struct {
	SpanName string
}

func New(name string) Instrumented {
	return Instrumented{
		SpanName: name,
	}
}

func (s *Instrumented) SpanTrace(ctx context.Context, name string, tags map[string]any) (context.Context, opentracing.Span) {
	return jaeger.SpanTrace(ctx, fmt.Sprintf("%s.%s", s.SpanName, name), tags)
}
