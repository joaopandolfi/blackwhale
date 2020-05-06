package mysql

import (
	"github.com/joaopandolfi/blackwhale/configurations"
	"github.com/joaopandolfi/blackwhale/remotes/sqldriver"
)

type pDriver struct {
	sqldriver.SqlDriver
	Case string
}

// Driver postgresqk
var Driver pDriver

// Init function used to initialize mysql database
func Init() {
	Driver = pDriver{Case: "lower"}
	Driver.Init("postgres", configurations.Configuration.PostgreSQL)
}
