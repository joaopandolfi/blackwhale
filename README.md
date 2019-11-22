# blackwhale
Go web Framework


# Main Example

```
package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joaopandolfi/blackwhale/configurations"
	"github.com/joaopandolfi/miia_api/routes"
	"github.com/unrolled/secure"

	"github.com/joaopandolfi/blackwhale/handlers"
	"github.com/joaopandolfi/blackwhale/remotes/mysql"
	"github.com/joaopandolfi/blackwhale/utils"
)

func configInit() {
	configurations.Load()
	mysql.Init()
	// Precompile html pages
	routes.Precompile()
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