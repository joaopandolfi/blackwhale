# Jaeger tracing implementation

## Config
You need setup those envs
```
    - JAEGER_SERVICE_NAME=service-x
    - JAEGER_AGENT_HOST=jaeger
    - JAEGER_SAMPLER_TYPE=const
    - JAEGER_SAMPLER_PARAM=1
    - JAEGER_REPORTER_LOG_SPANS=true
```

## Init
```GoLang
	tracer, closer := jaeger.Init(thisServiceName)
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)
```

## Start tracing
```GoLang
    newCtx, span := jaeger.SpanTrace(r.Context(), "test", map[string]interface{}{})
    defer span.Finish()
```

## Start tracing from http context
```GoLang
    newCtx, span := jaeger.StartSpanFromRequest(tracer, r)
    defer span.Finish()
```

## Running local jaeger host
```yaml
version: '3'
services:
  service-a:
    image: service-a
    ports:
      - "8081:8081"
    environment:
      - JAEGER_SERVICE_NAME=service-a
      - JAEGER_AGENT_HOST=jaeger
      - JAEGER_SAMPLER_TYPE=const
      - JAEGER_SAMPLER_PARAM=1
      - JAEGER_REPORTER_LOG_SPANS=true
  jaeger:
    image: jaegertracing/all-in-one
    ports:
      - "16686:16686"
```
