# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [2.0.0] - 2026-09-26

New major line (module path `github.com/joaopandolfi/blackwhale/v2`). Breaking release.

### Added

- MIT `LICENSE` (root) and upstream attribution for the vendored `handlers/conjson`.
- `Makefile` (build/vet/test/lint/coverage/integration) and GitHub Actions CI
  (unit, lint with legacy baseline, integration behind the `integration` tag + compose).
- `docker-compose.integration.yml` for integration services.
- `remotes/jaeger`: OpenTelemetry tracing over OTLP HTTP (`OTEL_EXPORTER_OTLP_ENDPOINT`,
  OTel samplers). `jaeger.Init(service)` now actually uses the `service` parameter and
  sets the global tracer provider.
- `remotes/mongo/v2`: methods accept `context.Context`.

### Changed

- Module path renamed to `.../v2`; Go `1.26`.
- `interface{}` → `any`, `chan bool` → `chan struct{}`, `ioutil` → `io`/`os`,
  `pkg/errors` → stdlib `fmt.Errorf` with `%w`.
- Dependencies: `golang-jwt` v4 → v5, `go-redis` v8 → v9, `go-playground/validator`
  v9 → v10, `streadway/amqp` → `rabbitmq/amqp091-go` (package `amqp091`),
  `mongo-driver` v1 → v2, `gorilla/mux` → `go-chi/chi` v5.
- Tracing: `opentracing` + `uber/jaeger-client-go` removed in favor of OpenTelemetry
  (`span.Finish()` → `span.End()`; consumers drop `opentracing.SetGlobalTracer`).
- Router: public APIs `handlers.HandleTokenPermissions`, `handlers.QuietHandleTokenPermissions`
  and `prometheus.InjectMiddleware` now take `*chi.Mux`; `handlers.GetVars` reads
  chi route params; the prometheus middleware labels by chi route pattern.
- `remotes/mongo`: rewritten on the official `mongo-driver/v2`. `Session` wraps a shared
  `*mongo.Client` per URL; `GetCollection`/`GetCollectionOnDB` return `*mongo.Collection`;
  `GetNextID`/`CreateIndex`/`Run` reimplemented on the driver API. TLS: `ssl=true`
  normalized to `tls=true` with `InsecureSkipVerify` (same behavior as old `NewSessionSsl`).
- `utils` JWT functions are now wrappers around `remotes/jwt` (canonical
  implementation); `utils.Token` is an alias of `remotes/jwt.Token` (adds `Broker`).
- Config secrets (`MYSQL_USER`, `MYSQL_PASSWORD`, `RESET_HASH`, `BCRYPT_SECRET`,
  `SESSION_SECRET`) are read from environment instead of hardcoded defaults;
  the expired Hasura JWT in tests is gone (`HASURA_SYSTEM_TOKEN` required).

### Removed

- `gopkg.in/mgo.v2` (dead 2019 driver) and its unused API surface:
  `mongo.NewService`, `mongo.GenericInsert`, `mongo.NewSessionSsl`,
  `mongo.NewSessionSSLMETHOD2`, `mongo.NewSessionManual`, `mongo.NewCustomSession`,
  `mongo.FlushPull`.
- `opentracing`-typed helpers (`jaeger.Extract`); `jaeger.Inject` now takes only the
  request (context carries the span).
- `InsecureSkipVerify` on the generic `remotes/request` HTTPS client.

### Fixed

- `models/cryptable`: AES key captured at package init (before `configurations.Load()`)
  — now read at call time, so encryption no longer silently uses a stale/empty key.
  `SetAesKey` remains as an explicit override.
- `cron`: ephemeral jobs could never fire (ticker drained only on a coincidental tick);
  the pending tick is now flushed before evaluation.
- `remotes/cache`: the lazy-init 40s timeout panicked in a goroutine (process crash);
  it now keeps waiting with a periodic error log. `waitListenners` protected by a mutex
  (data race). `SafeCache` methods guard against a not-yet-injected cache instead of
  nil-dereferencing.
- `remotes/graphite`: IPv6 hosts no longer break the dial (`net.JoinHostPort`).
- `remotes/hasura` tests no longer embed an expired JWT.
- `remotes/sqldriver`: `MaxSizeMbUpload` overflow (`10<<55` → `10<<20`) and case
  normalization preserved (`toCase`); `gosqljson` API aligned with the current version.
- `remotes/jwt` and `utils/crypt`: JSON struct tags on `Token` fixed (`json:'id'` → `json:"id"`).
- Test hygiene: flaky cron test and jaeger panic eliminated; integration tests gated
  behind the `integration` build tag.

## [1.9.0] - 2026-09-26

Last release of the v1 line (branch `modernization`): hygiene, CI/lint/compose tooling,
security fixes (secrets to env, MIT license, conjson attribution), build/vet/test fixes.
Non-breaking relative to v1.8.3.
