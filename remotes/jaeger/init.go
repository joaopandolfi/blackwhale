package jaeger

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
	"go.opentelemetry.io/otel/trace"
)

// Init returns a tracer provider that exports spans to Jaeger via OTLP.
// It also sets the provider as the global one, so SpanTrace and
// StartSpanFromRequest work without an explicit tracer argument.
// OTEL_EXPORTER_OTLP_ENDPOINT overrides the default http://localhost:4318.
func Init(service string) (trace.TracerProvider, io.Closer) {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:4318"
	}

	opts := []otlptracehttp.Option{otlptracehttp.WithEndpointURL(endpoint)}
	if !strings.HasPrefix(endpoint, "https://") {
		opts = append(opts, otlptracehttp.WithInsecure())
	}
	exporter, err := otlptracehttp.New(context.Background(), opts...)
	if err != nil {
		panic(fmt.Sprintf("cannot init jaeger exporter: %v", err))
	}

	resource, err := sdkresource.Merge(
		sdkresource.Default(),
		sdkresource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceName(service)),
	)
	if err != nil {
		panic(fmt.Sprintf("cannot init jaeger resource: %v", err))
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource),
	)
	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return provider, &providerCloser{provider}
}

type providerCloser struct {
	provider *sdktrace.TracerProvider
}

func (c *providerCloser) Close() error {
	return c.provider.Shutdown(context.Background())
}
