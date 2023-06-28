package instrumentable

import (
	"context"
	"fmt"

	"github.com/joaopandolfi/blackwhale/remotes/jaeger"
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

func (s *Instrumented) SpanTrace(ctx context.Context, name string, tags map[string]interface{}) (context.Context, opentracing.Span) {
	return jaeger.SpanTrace(ctx, fmt.Sprintf("%s.%s", s.SpanName, name), tags)
}
