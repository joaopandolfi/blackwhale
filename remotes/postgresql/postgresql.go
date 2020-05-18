package postgresql

import (
	"github.com/joaopandolfi/blackwhale/configurations"
	"github.com/joaopandolfi/blackwhale/remotes/sqldriver"
)

type pDriver struct {
	sqldriver.SqlDriver
	Case string
}

// Driver postgresql
var Driver pDriver

// Pool postgresql
var Pool map[string]pDriver

// Init function used to initialize mysql database
func Init() {
	Driver = pDriver{Case: "lower"}
	Driver.Init("postgres", configurations.Configuration.PostgreSQL)
}

// MakePool driver
func MakePool(key, url string) {
	d := pDriver{Case: "lower"}
	d.Init("postgres", url)

	if Pool == nil {
		Pool = make(map[string]pDriver)
	}
	Pool[key] = d
}

// ClosePool database
func ClosePool() {
	if Pool == nil {
		return
	}
	for _, p := range Pool {
		p.Close()
	}
}
