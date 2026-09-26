# blackwhale — guide for AI agents

Base library (not a framework) for Go web services: config, HTTP handlers/middlewares,
JWT auth, cache, cron, and remote drivers. Module `github.com/joaopandolfi/blackwhale/v2`,
Go 1.26, MIT.

## Commands

```sh
make build      # go build ./...
make vet        # go vet ./...
make test       # go test -count=1 ./...   (unit, no infra needed — must stay green)
make lint       # golangci-lint (config: .golangci.yml; only-new-issues vs legacy baseline)
make integration # needs docker-compose.integration.yml up (see README)
```

Integration tests are behind the `integration` build tag:
`go test -tags=integration ./remotes/mongo/v2/... ./remotes/rabbitmq/... ./remotes/beanstalkd/...`
Compile-check them any time with `go vet -tags=integration ./...`.

## Architecture rules (v2)

- **Public API is `net/http`-shaped.** Handlers are `http.Handler`/`http.HandlerFunc`;
  middlewares are `func(http.HandlerFunc) http.HandlerFunc` (`handlers.Middleware`).
  Routers are `*chi.Mux` (chi v5) — never gorilla types in signatures.
- **`configurations.Configuration` is the global config singleton** (13/18 consumer
  projects depend on it). Do not invert it without a v3 decision. Config is loaded once
  via `configurations.Load()`; secrets come from env, never from code.
- **Dependency direction**: `remotes/* → configurations` (allowed),
  `utils → handlers/conjson` (allowed, one-way — conjson is a vendored MIT codec).
  Never import upward: `handlers` may not be imported by `configurations`/`utils`
  (except conjson), and `remotes` must not import `handlers`.
- **Init functions are fail-fast and panic** on misconfiguration
  (`configurations.Load`, `jaeger.Init`, `mongo.GetSession`, ...). Consumers wrap
  `main` with a `recover()` ("resilient" pattern). Keep this contract; return errors
  only for call-time/usage failures.
- **Context**: `remotes/mongo/v2` methods take `context.Context`. The legacy
  `remotes/mongo` facade (consumer contract) uses `context.Background()` internally —
  do not change its signatures without a v3.

## Conventions

- `any`, not `interface{}`; `fmt.Errorf("...: %w", err)` with wrapped errors;
  no `pkg/errors`, no `ioutil`, no `panic` for recoverable errors.
- No comments/docstrings unless the behavior is not inferable from the signature
  (a pre-commit hook enforces this). Docstrings start with the symbol name.
- Commits: plain English subject, one logical change per commit, footer:
  `Ultraworked with [Sisyphus](https://github.com/code-yeongyu/oh-my-openagent)` +
  `Co-authored-by: Sisyphus <clio-agent@sisyphuslabs.ai>`.
- `gofmt` clean (`gofmt -l .` must print nothing) — the hook and CI check it.

## Gotchas

- **Consumers**: 18 projects in `~/projetos` import this lib. Any signature change to
  core packages (`utils`, `configurations`, `handlers`, `remotes/mongo`, `remotes/request`)
  is a breaking change for them — check usage before changing signatures.
- **JWT**: canonical implementation is `remotes/jwt` (explicit secret). `utils` JWT
  functions are config-secret wrappers; `utils.Token` is an alias of `jwt.Token`.
  Do not duplicate JWT logic.
- **Mongo**: `mgo.v2` is dead. `remotes/mongo` (facade) and `remotes/mongo/v2` (thin
  client) both use `mongo-driver/v2`. Consumer code must use the official driver API
  (`FindOne`/`InsertOne`/`UpdateOne`/`DeleteMany`) and `mongo-driver/v2/bson`.
- **Tracing**: OpenTelemetry over OTLP HTTP (default `http://localhost:4318`,
  `OTEL_EXPORTER_OTLP_ENDPOINT` overrides). `jaeger.Init(service)` sets the global
  provider — code calls `jaeger.SpanTrace`/`StartSpanFromRequest` with no tracer arg.
  Spans end with `span.End()`.
- **`utils/snake_case`** imports `handlers/conjson` — intentional one-way flow; do not
  "fix" it into a cycle.
- **gosqljson** has no tagged versions (pseudo-version pinned in go.mod); its API is
  `QueryToMaps`/`QueryToArrays` + `AsIs`/`Lower`/`Upper`/`Camel`.
- **OTel semconv** must match the SDK schema URL (sdk v1.46.0 → `semconv/v1.43.0`) or
  `resource.Merge` panics.
- **chi v5.3.x**: `Methods(...)` does not exist on `chi.Router`; use
  `*chi.Mux.MethodFunc(method, pattern, fn)` per method. Route params:
  `chi.RouteContext(r.Context()).URLParams`.
- **mongo-driver v2**: `mongo.Connect` takes no context; `RunCommand` returns a lazy
  `*SingleResult` (must `Decode` to execute) and needs order-preserving commands
  (`bson.D`, not `bson.M`).
- **`handlers/conjson` is vendored** (upstream Rican7/conjson, MIT — see its LICENSE).
  Do not rewrite; attribute changes upstream.

## Status

- `master` = v1.8.3 (frozen). `modernization` = v1 line, tag `v1.9.0` (frozen).
  `v2` = current development (this guide), phases 0–6 of `docs/MODERNIZATION.md` done.
- Open items: rotate any exposed secrets (maintainer), push branches when authorized.
