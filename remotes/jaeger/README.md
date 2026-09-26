# Jaeger tracing (OpenTelemetry)

Spans are exported to Jaeger over OTLP. Jaeger v1.35+/v2 natively accepts OTLP.

## Config

| Env | Default | Purpose |
|---|---|---|
| `OTEL_EXPORTER_OTLP_ENDPOINT` | `http://localhost:4318` | OTLP HTTP endpoint |
| `OTEL_TRACES_SAMPLER` | `parentbased_always_on` | Standard OTel sampler (`always_on`, `traceidratio`, ...) |
| `OTEL_TRACES_SAMPLER_ARG` | - | Sampler argument (e.g. ratio `0.1`) |

## Init

```go
provider, closer := jaeger.Init(thisServiceName)
defer closer.Close()
```

`Init` sets the provider as the global one, so the span helpers below work without
an explicit tracer argument.

## Start tracing

```go
ctx, span := jaeger.SpanTrace(r.Context(), "test", map[string]any{"k": "v"})
defer span.End()
```

## Start tracing from an inbound HTTP request

```go
ctx, span := jaeger.StartSpanFromRequest(r, "op.name")
defer span.End()
```

## Propagate to an outbound HTTP request

```go
req = req.WithContext(ctx)
jaeger.Inject(req)
```

## Running a local Jaeger host

```yaml
version: '3'
services:
  service-a:
    image: service-a
    ports:
      - "8081:8081"
    environment:
      - OTEL_EXPORTER_OTLP_ENDPOINT=http://jaeger:4318
      - OTEL_TRACES_SAMPLER=always_on
  jaeger:
    image: jaegertracing/all-in-one
    ports:
      - "16686:16686"
```
