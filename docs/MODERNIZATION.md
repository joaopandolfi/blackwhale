# blackwhale — Análise & Plano de Modernização

**Data:** 2026-09-26 · **Status:** aprovado para execução (fases 0–2)
**Repo:** `github.com/joaopandolfi/blackwhale` — base lib Go (framework web + remotes + models + helpers)
**Método:** análise por 4 agentes paralelos (estrutura, docs, QA/CI, padrões desatualizados) + digest de best-practices 2025–2026.

---

## 0. Veredito

A hipótese ("muita coisa desatualizada + documentação ruim") está **confirmada e é pior do que parecia**. O repo é uma base lib de ~2020–2021 (primeiro commit 2019-11, último 2024-11-26) cuja *diretiva* do Go foi atualizada em 2024, mas cujo código, dependências, testes e documentação seguem congelados naquela época. **Não existe nenhuma infraestrutura de QA** (CI, linter, `.gitignore`, LICENSE). E há **problemas de segurança reais** (segredos hardcoded) além de **um bug latente**.

**Status atual (2026-09-26, pós-Fase 0):** `go build ./...`, `go vet ./...` e a compilação de todos os testes estão **verdes**. Concluído nesta rodada:
1. Regressão de segredos do WIP revertida (§1.1);
2. **Build break existente em HEAD corrigido** — `remotes/sqldriver` chamava `gosqljson.QueryDbTo*`, API que **não existe** na única versão do pacote (pseudo-version 2023-04, repo sem tags); adaptado para `QueryToMaps`/`QueryToArrays` mantendo a API pública (`theCase string`);
3. Struct tags inválidos `json:'id'` (aspas simples → claims JWT serializados com nome errado) corrigidos em `remotes/jwt/entities.go` e `utils/crypt.go` (flag do `go vet`);
4. `net.Dial` com formato IPv6 quebrado em `remotes/graphite/driver.go` → `net.JoinHostPort`.

Infra de QA criada em sessão paralela (ainda untracked): `.github/workflows/ci.yml` (3 jobs: unit/lint/integration), `.golangci.yml` (v2, conservador), `docker-compose.integration.yml`, `.gitignore`.

---

## 1. ⚠️ Descoberta crítica: já existe WIP no working tree

Antes de qualquer ação, o working tree **não é lixo** — contém uma **Fase 1 parcial já em curso** (provavelmente de uma sessão anterior). Separar o que está *staged* do que está *unstaged*:

**Staged (já no índice):**
- `.DS_Store` → **removido** (untrack já feito)
- `remotes/mongo/test.go` → **removido** (o arquivo morto com `func main()` em package de lib)
- `handlers/request_test.go` → teste novo adicionado (61 linhas; versão staged **sem** build tag)
- `models/cryptable/encryptable.go` → removida a tag `json:"Crypted"` redundante (no-op p/ JSON)

**Unstaged (só no working tree):**
- `README.md` → nova seção `# Testing` (unit + integration atrás da tag, `docker-compose.integration.yml`, overrides por env var)
- `configurations/configurations.go` → **MUDANÇA PERIGOSA (ver §2)** + fixes de typo (`locahost`→`localhost`, `localhosy`→`localhost`) + fix real de bug (`MaxSizeMbUpload: 10<<55` overflow → `10<<20` = 10 MB)
- `handlers/request_test.go` → adiciona `//go:build integration` (corrige a versão staged)
- `remotes/{beanstalkd,graphite,hasura,mongo/v2,rabbitmq}/*_test.go` → **adicionam `//go:build integration`** e **param endpoints por env var** (`MONGO_HOST`, `RABBITMQ_URL`, `BEANSTALKD_URL`, `GRAPHITE_HOST/PORT`, `HASURA_URL/SYSTEM_TOKEN`), trocando os IPs internos hardcoded (`10.0.0.2`, `10.0.0.14`) por `127.0.0.1`
- `remotes/request/https.go` → **removido `InsecureSkipVerify: true`** (fix de segurança) + `fmt.Println` → `log.Printf`

**Conclusão:** esse WIP já cobre boa parte da Fase 1 do plano (tags de integration + endpoints por env var + higiene). A decisão correta é **revisar e completar**, **não descartar**. Falta: o `docker-compose.integration.yml` referenciado pelo README **não existe ainda**, e há a regressão de segurança abaixo.

### 1.1 ✅ Regressão de segurança no WIP — revertida (estava pendente antes de qualquer commit)

Em `configurations/configurations.go`, o WIP **substituiu segredos dirigidos por config por valores hardcoded**:

| Local | Antes (correto) | Depois (WIP — errado) |
|---|---|---|
| `LoadFromMap` | `JWTSecret: fconf["JWT_SECRET"]` | `JWTSecret: "supersecretkey"` |
| `LoadFromMap` | `AESKEY: fconf["AES_KEY"]` | `AESKEY: "0123456789abcdef0123456789abcdef"` |
| `Load()` | `JWTSecret: ""` | `JWTSecret: "supersecretkey"` |
| `Load()` | `AESKEY: "-weak key :( -"` | `AESKEY: "0123456789abcdef0123456789abcdef"` |

Isso **quebra o propósito de segredo configurável** e é exatamente o anti-padrão que a análise flagou (segredos hardcoded em `configurations.go`). **Recomendação: reverter essas 4 linhas para leitura de config/env** (ou vazio + obrigatório em runtime) antes de qualquer commit. Os demais changes do arquivo (typos, `MaxSizeMbUpload`, `Debug`) podem ser mantidos.

**Resolvido nesta rodada:** `JWTSecret`/`AESKEY` voltaram a `fconf["JWT_SECRET"]`/`fconf["AES_KEY"]` em `LoadFromMap` e `""` em `Load()`; typos, `MaxSizeMbUpload` e `Debug` mantidos. Os demais segredos default em `Load()` (MySQL `root`, securecookie etc.) seguem pendentes na Fase 2.

---

## 2. Estado atual (números)

| Métrica | Valor |
|---|---|
| Escala | 84 arquivos `.go`, ~5.900 linhas de código, ~770 de testes |
| Pacotes | 7 áreas: `configurations`, `cron`, `handlers` (3 sub), `middlewares`, `models` (5 sub), `remotes` (16 sub!), `utils` (4 sub) |
| Atividade | Primeiro commit 2019-11, último 2024-11-26 (`fix gzip request`) |
| go.mod | `go 1.23.2` — **EOL desde ago/2025** (suportadas hoje: 1.26 e 1.27; toolchain local 1.26.5) |
| README | 2021, 73 linhas, exemplo único **não compila** (`handlers.Response` mudou de assinatura) |
| Testes | 27 funções em 16 arquivos; ~1/3 exigia serviço vivo em IP interno hardcoded |
| CI / lint / higiene | **zero** — sem workflow, Makefile, `.golangci.yml`, `.gitignore`, LICENSE, Dockerfile |

> Nota: `models` e `remotes` **nem são pacotes Go** (só diretórios-namespaço, 0 `.go` no topo) — nem poderiam ter package comment.

---

## 3. Findings por área

### 3.1 Dependências desatualizadas (go.mod)

| Dependência | Situação |
|---|---|
| `gorilla/mux` v1.8.1 | **maintenance mode** (arquivado ~2022) |
| `golang-jwt/jwt/v4` | v5 é a atual |
| `gopkg.in/mgo.v2` | **descontinuado desde 2019**, convivendo com `mongo-driver v1` (v2 atual) — 2 drivers Mongo no repo |
| `go-redis/v8` | v9 atual |
| `go-playground/validator v9` | `+incompatible`; v10/v11 atuais |
| `pkg/errors` | superseded por stdlib `errors` + `%w` |
| `opentracing` + `uber/jaeger-client-go` | **superseded por OpenTelemetry** (`otel v1.29` já é dep indireta) |
| `streadway/amqp` | superseded por `rabbitmq/amqp091-go` |
| `flosch/pongo2` (2020), `kr/beanstalk` (2018) | antigos, sem movimento |

### 3.2 Arquitetura

- **Singletons mutáveis onipresentes** — 46 `var` package-level em 21 arquivos: `configurations.Configuration`, `cron.Get()`, `cache.Get()`, `mysql.Driver`, `postgresql.Pool`, estado global do graphite, `rabbitmq.conn`…
- **Bug latente real:** `models/cryptable/encryptable.go:10` — `var aesKey = configurations.Configuration.Security.AESKEY` copia o valor **no init do package, antes de `Load()` rodar** → chave vazia a menos que algo re-copie.
- **Duplicações:** JWT implementado 2× (~87 linhas idênticas em `remotes/jwt` e `utils/crypt.go`); Mongo com 2 gerações; Graphite com 2 APIs (v1 `Driver` + v2 `Graphite`) no mesmo package.
- **Ciclos de grupo:** `configurations ↔ remotes` e `handlers ↔ utils` (via `utils/snake_case → handlers/conjson`). Não quebram o build, mas travam refator.
- **Código de terceiro vendado:** `handlers/conjson` é cópia do `Rican7/conjson` (copyright intacto, README com badges do projeto original — e **sem LICENSE** no repo).
- **Código morto:** `remotes/mongo/test.go` (já staged p/ remoção) + 6 blocos `//panic(err)` comentados em `sqldriver`.
- **13 `panic(` ativos em código de biblioteca** (mongo, pubsub, jaeger, cache, configurations, dao).
- **Idioms modernos não adotados:** `interface{}` 164× / `any` 0×; `chan bool` 18×; `ioutil` em 2 arquivos.
- **`context` descartado:** `remotes/mongo/v2/client.go:27,44,75,77,89` substitui o ctx recebido por `context.TODO()`.

### 3.3 Segurança ⚠️

- **Segredos hardcoded em `configurations/configurations.go`:** `JWTSecret` default (:198,:291), **chave AES-256 de 32 bytes hardcoded** (:199,:292), credenciais MySQL `root` em `Load()` (:237-242), token `ResetHash` (:258), salt `BCryptSecret` (:267), chave securecookie de 64 bytes (:272). Se o repo foi público/compartilhado, **trate esses valores como comprometidos e rotacione**.
- JWT hardcoded em teste (`remotes/hasura`, expirou em 2022) e IPs internos (`10.0.0.2`, `10.0.0.14`) em testes. *(O WIP no working tree já endereça os IPs → env var; o JWT do hasura segue como fallback default — ver §1.)*
- `.DS_Store` tracked + sem `.gitignore`. *(WIP já untrackeia `.DS_Store`; falta o `.gitignore`.)*
- `InsecureSkipVerify: true` em `remotes/request/https.go`. *(WIP já remove — direção correta, mas é change de comportamento p/ consumidores.)*

### 3.4 Documentação (confirma "está ruim")

- **Zero `doc.go` no repo.** As únicas 2 package comments de todo o repo estão no código vendado (conjson).
- README de 2021: diz "Go web Framework" (é uma base lib com 16 remotes); exemplo usa `handlers.Response(w, "HELLOOU")` mas a assinatura real é `Response(w, resp, status)` — **não compila**; fence quebrado (``` ```package main ```); referência morta a `fvbock/endless`.
- Cobertura de doc comments em exported symbols: **cron 0%**, **middlewares 0%**, configurations ~33%, utils ~31% — a API de logging que o próprio README mostra (`Info`, `CriticalError`) está **100% sem doc**.
- Sem LICENSE, CHANGELOG, examples/, CONTRIBUTING.

### 3.5 Testes e CI (a infraestrutura que não existe)

- **Nenhum CI, nenhum linter, nenhum Makefile, nem Dockerfile.**
- 6 packages de teste exigiam infra viva **sem build tag** (mongo/v2, rabbitmq, beanstalkd, graphite, hasura, handlers) → `go test ./...` falhava em máquina limpa. *(WIP já adiciona as tags.)*
- 30s+ de `time.Sleep` em testes "unitários"; mocks (`MockJob`, `client_mock.go`, `pubsub/mock.go`) **shipados como arquivos de produção**, não de teste.
- `configurations` e `middlewares`: zero testes.

---

## 4. Plano de ajuste

Sequenciado p/ reduzir risco: primeiro o que habilita o resto, depois segurança, depois o que é breaking. **Fases 0–1 já estão ~50% feitas pelo WIP no working tree** (ver §1).

### Fase 0 — Higiene e pré-condições — **~90% done**
1. ✅ `.DS_Store` untracked (já staged no WIP)
2. ✅ `remotes/mongo/test.go` removido (já staged no WIP)
3. ✅ `.gitignore` criado (sessão paralela; cobre `.DS_Store`, `.yolo.json`, `coverage.out`, `go.work*`, `vendor/`, IDE, macOS)
4. ✅ **WIP revisado**: regressão de segredos revertida (§1.1); mantidos typos + `MaxSizeMbUpload` + `Debug` + tags de integration + env vars
5. ✅ `docker-compose.integration.yml` criado (sessão paralela: mongo:7, rabbitmq:3, beanstalkd c/ healthchecks)
6. ✅ **Consumidores levantados** (§5 — 18 projetos)
7. ✅ **Build break de HEAD corrigido** (`remotes/sqldriver` × gosqljson) + struct tags inválidas (vet) + dial IPv6 (graphite) — `go build`/`go vet`/test-compile verdes
8. 🔲 Decisão: commitar o WIP revisado (só com pedido explícito)

### Fase 1 — Infraestrutura de QA — **~70% done**
1. ✅ Tags `//go:build integration` + endpoints por env var (WIP)
2. 🔲 `go.mod`: `go 1.26` (melhor prática p/ lib = N-1; hoje suporta-se 1.26/1.27) + `toolchain`
3. ✅ `.golangci.yml` (formato v2, conservador: govet, ineffassign, staticcheck, unused + gofmt) — sessão paralela · 🔲 `Makefile` (build/vet/test/lint/coverage)
4. ✅ GitHub Actions (`.github/workflows/ci.yml`, sessão paralela): 3 jobs — unit (vet+build+test), lint (golangci-lint v2, only-new-issues), integration (compose up + testes tagados)
5. 🔲 **Primeiro tag semântico** (`v0.1.0` ou `v1.0.0`) — a partir daqui consumidores conseguem pinar

### Fase 2 — Segurança (~1–2 dias, independente)
1. 🔲 Tirar **todos** os defaults de segredo de `configurations.go` (carregar só de arquivo/Vault/env)
2. 🔲 **LICENSE** (obrigatório — há código vendado com copyright e nenhum license no repo, risco legal)
3. 🔲 Resolver conjson: adicionar license upstream OU virar dependência real
4. 🔲 Remover JWT expirado do teste do hasura (ou exigir env var obrigatória)
5. 🔲 Rotacionar o que foi exposto (se o repo já foi público/compartilhado)

### Fase 3 — Deps e código (mecânico → arriscado)
1. 🔲 **Mecânico e seguro** (PR separado): `interface{}`→`any` (164), `chan bool`→`chan struct{}` (18), `ioutil`→`io`, `pkg/errors`→stdlib, `errors.New(fmt.Sprintf)` (4 locais), 26 `fmt.Errorf` sem `%w`
2. 🔲 **Upgrades c/ risco de regressão** (1 por PR): jwt v4→v5 (c/ dedup), redis v8→v9, mongo-driver v1→v2, validator v9→v10, `streadway/amqp`→`rabbitmq/amqp091-go`, opentracing/jaeger→**OpenTelemetry**
3. 🔲 **Matar `mgo.v2`:** migrar `remotes/mongo` pro driver oficial e apagar a geração antiga
4. 🔲 **gorilla/mux → chi v5** — o maior item. Regra de ouro: API pública continua **shaped em `net/http`** (handlers = `http.Handler`, middlewares = `func(http.Handler) http.Handler`)

### Fase 4 — Arquitetura (o caro; depende da decisão de versão)
1. 🔲 Quebrar ciclos: `configurations↔remotes` (vault) e `handlers↔utils` (snake_case → package próprio)
2. 🔲 Estratégia de singletons: manter só `config` como singleton raiz; demais remotes → construtores injetáveis c/ default
3. 🔲 Corrigir o bug do `aesKey` (cryptable)
4. 🔲 Dissolver `utils`: logging → `log/slog`; demais helpers p/ packages nomeados ou `internal/`
5. 🔲 Dedup: JWT 2→1, Graphite 2 APIs→1, Mongo 2→1
6. 🔲 13 `panic` → erros retornados; `context.Context` nos funcs de I/O

### Fase 5 — Documentação (~1 semana, paralelizável)
1. 🔲 README reescrito: o que é (base lib, não framework), install (`go get` + `GOPRIVATE`), quickstart que **compila**, tabela de pacotes, versões Go suportadas, política de versionamento
2. 🔲 `doc.go` + package comment em todo package público; doc comments começando por `utils` (logging) e `cron`
3. 🔲 `Example*` com `// Output:` nos packages públicos (dobram como teste de regressão)
4. 🔲 CHANGELOG (Keep a Changelog); atualizar README stale do jaeger (assinatura desatualizada) e do pubsub (zero API)

### Fase 6 — Testes (~1–2 semanas, paralelizável)
1. 🔲 Unit tests de `configurations` (parse/load) e `middlewares` (gzip) — lógica pura, sem infra
2. 🔲 Testes de HTTP real → `httptest`; eliminar sleeps (`Eventually`/canais); mocks p/ `_test.go`/`internal`
3. 🔲 Meta: suite unitária 100% verde em CI sem infra; integration atrás da tag + compose

---

## 5. Consumidores

**18 projetos** em `~/projetos` importam `github.com/joaopandolfi/blackwhale`:

| Área | Projetos |
|---|---|
| `tools/` | `remote_config/src`, `PoW-Shield-Go/server` |
| `w3care/tools/` | `ms_payment`, `service-watcher/src`, `ms_logger/src`, `babel`, `messenger/tools/bomber`, `messenger/go_broker`, `blank-golang-project/echo`, `blank-golang-project/mux`, `smeagle/src`, `drug_interaction/app`, `oraculo/src`, `ms_file_cript/go_src` |
| `surgical/` | `backend/src` |
| `gent/` | `sup-bot/src`, `code/whatsapp-integrator-backend/src` |

Nota p/ Fase 3 (sqldriver): o único `theCase` observado em todo o corpus (blackwhale + consumidores) é `"lower"` — o mapper `toCase` cobre `lower/upper/camel` e o default é `AsIs`, então nenhum consumidor muda de comportamento.

**Ainda falta mapear** (pré-requisito p/ fases 3–4): qual API pública cada um usa (handlers? middlewares? quais remotes?) — dimensiona o esforço de migração e a janela de coexistência do v2.

---

## 6. A decisão que muda tudo: v1 ou v2?

As fases 3–4 (API `net/http`-shaped, `context`, chi, dedup) **são breaking por natureza**. Duas rotas:

1. **Corte v2 (recomendado):** fases 0–2 saem como `v1.x` (non-breaking); fases 3–6 viram o `v2.0.0` com module path `/v2` — o rename de import em **um** momento só, com janela de coexistência e `// Deprecated:` no v1. Última chance barata do rename.
2. **Manter v1 (estilo testify):** mudar APIs c/ mais cuidado, sem rename, aceitando API mista por um tempo. Menos dor de migração, mais dívida acumulada.

Impacta todos os services consumidores — por isso o levantamento de consumidores (§5) vem antes.

---

## 7. Próximos passos imediatos

1. ✅ Reverter a regressão de segredos em `configurations.go` (§1.1)
2. ✅ Criar `.gitignore` e `docker-compose.integration.yml` (sessão paralela)
3. ✅ Confirmar build (`go build ./...`), vet e compilação dos testes — **verdes**
4. ✅ Preencher §5 (consumidores — 18 projetos)
5. 🔲 Decidir: commitar o WIP revisado? (exige pedido explícito)
6. 🔲 Escolher a rota de versionamento (§6)
7. 🔲 Mapear a superfície de API usada por cada um dos 18 consumidores (pré-requisito p/ fases 3–4)
