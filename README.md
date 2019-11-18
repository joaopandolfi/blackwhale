# blackwhale
Go web Framework


# Main Example

```
func configInit() {
	configurations.Load()
	mysql.Init()
	// Precompile html pages
	routes.Precompile()
}

func resilient(){
	utils.Info("[SERVER] - Shutdown")

	if err := recover(); err != nil {
		utils.CriticalError("[SERVER] - Returning from the dark",err)
		main()
	}
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

	// Routes consist of a path and a handler function.
	routes.Register(r)

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