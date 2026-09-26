# blackwhale

Base library for Go web services: configuration, HTTP handlers/middlewares, auth (JWT),
caching, cron, and remote drivers (mongo, mysql, redis, rabbitmq, hasura, ...).
It is **not a framework** — it provides the pieces, you compose them.

- Module: `github.com/joaopandolfi/blackwhale/v2`
- Go: 1.26+
- License: MIT

## Install

```sh
go get github.com/joaopandolfi/blackwhale/v2
```

## Quickstart

```go
package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joaopandolfi/blackwhale/v2/configurations"
	"github.com/joaopandolfi/blackwhale/v2/handlers"
	"github.com/joaopandolfi/blackwhale/v2/remotes/mysql"
	"github.com/joaopandolfi/blackwhale/v2/utils"

	"github.com/unrolled/secure"
)

func configInit() {
	configurations.Load()
	mysql.Init()
}

func resilient() {
	utils.Info("[SERVER] - Shutdown")

	if err := recover(); err != nil {
		utils.CriticalError("[SERVER] - Returning from the dark", err)
		main()
	}
}

func relou(w http.ResponseWriter, r *http.Request) {
	handlers.Response(w, "HELLOOU")
}

func main() {
	defer resilient()

	//Init
	configInit()

	// Initialize chi Router
	r := chi.NewRouter()

	// Security
	secureMiddleware := secure.New(configurations.Configuration.Security.Options)
	r.Use(secureMiddleware.Handler)

	// Add routes
	r.Get("/", relou)

	// Bind to a port and pass our router in
	utils.Info("MI server listenning on", configurations.Configuration.Port)
	srv := &http.Server{
		Handler:      r,
		Addr:         configurations.Configuration.Port,
		WriteTimeout: configurations.Configuration.Timeout.Write,
		ReadTimeout:  configurations.Configuration.Timeout.Read,
	}

	err := srv.ListenAndServe()
	if err != nil {
		utils.CriticalError("Fatal server error", err.Error())
	}
}
```

Notes:

- `configurations.Load()` reads `config.yml` (or `$CONFIG_FILE`) plus env overrides and
  fills the global `configurations.Configuration`. Secrets come from env:
  `MYSQL_USER`, `MYSQL_PASSWORD`, `RESET_HASH`, `BCRYPT_SECRET`, `SESSION_SECRET`.
- Library init functions are fail-fast: they **panic** on misconfiguration. The
  `resilient()` recover wrapper above is the expected pattern for `main`.

## Packages

| Package | Purpose |
|---|---|
| `configurations` | Config load (file + env) and the global `Configuration` singleton |
| `handlers` | HTTP responses, request params, JWT/permission middleware (chi) |
| `handlers/conjson` | JSON codec with conventional-key transforms (vendored, MIT) |
| `handlers/errors` | Typed error responses |
| `middlewares` | HTTP middlewares (gzip) |
| `utils` | Logging, bcrypt, JWT helpers (wraps `remotes/jwt`), misc |
| `cron` | Cron jobs (incl. ephemeral jobs) |
| `models/compressible` | Transparent gzip-able fields |
| `models/cryptable` | AES field encryption (`Encryptable`) |
| `models/dao` | GORM-based DAO helpers |
| `models/permissions` | Permission check primitives |
| `models/transformers` | Field transformers |
| `remotes/beanstalkd` | Beanstalkd client |
| `remotes/cache` | Cache interface + redis/memory drivers, `SafeCache[T]` |
| `remotes/graphite` | Graphite metrics |
| `remotes/hasura` | Hasura GraphQL client |
| `remotes/instrumentable` | Span-name helper for tracing |
| `remotes/jaeger` | OpenTelemetry tracing (OTLP → Jaeger) |
| `remotes/jwt` | JWT HS256 issue/verify (canonical) |
| `remotes/mongo` | Mongo facade (singleton/pool, official driver v2) |
| `remotes/mongo/v2` | Thin mongo client with `context.Context` |
| `remotes/mysql` | MySQL (database/sql) |
| `remotes/postgresql` | PostgreSQL (database/sql) |
| `remotes/prometheus` | Prometheus metrics + HTTP middleware |
| `remotes/pubsub` | AMQP pub/sub (rabbitmq/amqp091-go) |
| `remotes/rabbitmq` | RabbitMQ driver |
| `remotes/request` | HTTP client (with retries) |
| `remotes/sqldriver` | SQL query helper (gosqljson) |

## Migrating from v1

v2 is a breaking release (module path `/v2`). Main changes:

- **Imports**: add `/v2` to every `github.com/joaopandolfi/blackwhale/...` import.
- **Router**: gorilla/mux → chi. `handlers.HandleTokenPermissions`,
  `handlers.QuietHandleTokenPermissions` and `prometheus.InjectMiddleware` take
  `*chi.Mux`; `handlers.GetVars` reads chi route params. Build routers with
  `chi.NewRouter()` (e.g. `r.Get("/x", h)` instead of `r.HandleFunc("/x", h).Methods("GET")`).
- **Mongo**: `mgo.v2` is gone. `remotes/mongo` now wraps the official
  `mongo-driver/v2`: sessions are shared `*mongo.Client`s per URL and
  `Session.GetCollection` returns `*mongo.Collection`. Migrate call sites to the
  official API (`FindOne`/`Find`, `InsertOne`, `UpdateOne`/`UpdateMany`,
  `DeleteMany`) and to `go.mongodb.org/mongo-driver/v2/bson` for `bson.M`.
- **Tracing**: opentracing/jaeger-client → OpenTelemetry. `jaeger.Init(service)`
  sets the global provider (drop `opentracing.SetGlobalTracer`), spans end with
  `span.End()`, config is OTel-native (`OTEL_EXPORTER_OTLP_ENDPOINT`, ...).
  See `remotes/jaeger/README.md`.
- **Secrets**: no hardcoded defaults. Set `MYSQL_USER`, `MYSQL_PASSWORD`,
  `RESET_HASH`, `BCRYPT_SECRET`, `SESSION_SECRET` in the environment.
- **AesKey fix**: `models/cryptable` now reads `AESKEY` from config at call time
  (v1 captured it at package init, before `configurations.Load()` — encryption
  silently used a stale key).
- `utils.Token` is an alias of `remotes/jwt.Token` (adds a `Broker` field).

## Testing

**Unit tests** (no external services needed):

```sh
go test ./...
```

**Integration tests** are gated behind the `integration` build tag and talk to
real services. For the services covered by `docker-compose.integration.yml`:

```sh
docker compose -f docker-compose.integration.yml up -d --wait
go test -tags=integration ./remotes/mongo/v2/... ./remotes/rabbitmq/... ./remotes/beanstalkd/...
docker compose -f docker-compose.integration.yml down -v
```

Integration endpoints can be overridden with environment variables:
`MONGO_HOST`, `MONGO_PORT`, `MONGO_DATABASE`, `RABBITMQ_URL`, `BEANSTALKD_URL`,
`GRAPHITE_HOST`, `GRAPHITE_PORT`, `HASURA_URL`, `HASURA_SYSTEM_TOKEN`.

`remotes/graphite` and `remotes/hasura` tests are integration-tagged but not
covered by the compose file — point them at your own instances (the Hasura one
expects a `patient` table in the metadata).

## Versioning

SemVer. The v1 line (module without `/v2`) is frozen at `v1.9.0` and kept for
coexistence; all development happens on v2.
