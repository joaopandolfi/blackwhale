# blackwhale
Go web Framework


# Main Example

```package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joaopandolfi/blackwhale/configurations"

	"github.com/unrolled/secure"

	"github.com/joaopandolfi/blackwhale/handlers"
	"github.com/joaopandolfi/blackwhale/remotes/mysql"
	"github.com/joaopandolfi/blackwhale/utils"
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

	// Initialize Mux Router
	r := mux.NewRouter()

	// Security
	secureMiddleware := secure.New(configurations.Configuration.Security.Options)
	r.Use(secureMiddleware.Handler)

	// Add routes
	r.HandleFunc("/", relou).Methods("GET")

	// Bind to a port and pass our router in
	utils.Info("MI server listenning on", configurations.Configuration.Port)
	srv := &http.Server{
		Handler:      r,
		Addr:         configurations.Configuration.Port,
		WriteTimeout: configurations.Configuration.Timeout.Write,
		ReadTimeout:  configurations.Configuration.Timeout.Read,
	}

	err := srv.ListenAndServe()
	//"github.com/fvbock/endless"
	///err := endless.ListenAndServeTLS("localhost:4242", "cert.pem", "key.pem", r)

	if err != nil {
		utils.CriticalError("Fatal server error", err.Error())
	}
}
```
# Testing

**Unit tests** (no external services needed):

    go test ./...

**Integration tests** are gated behind the `integration` build tag and talk to
real services. For the services covered by `docker-compose.integration.yml`:

    docker compose -f docker-compose.integration.yml up -d --wait
    go test -tags=integration ./remotes/mongo/v2/... ./remotes/rabbitmq/... ./remotes/beanstalkd/...
    docker compose -f docker-compose.integration.yml down -v

Integration endpoints can be overridden with environment variables:
`MONGO_HOST`, `MONGO_PORT`, `MONGO_DATABASE`, `RABBITMQ_URL`, `BEANSTALKD_URL`,
`GRAPHITE_HOST`, `GRAPHITE_PORT`, `HASURA_URL`, `HASURA_SYSTEM_TOKEN`.

`remotes/graphite` and `remotes/hasura` tests are integration-tagged but not
covered by the compose file — point them at your own instances (the Hasura one
expects a `patient` table in the metadata).
