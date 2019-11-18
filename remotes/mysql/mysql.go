package mysql

import (
	"github.com/joaopandolfi/blackwhale/configurations"
	"github.com/joaopandolfi/blackwhale/remotes/sqldriver"
)

type MySQLDriver struct {
	sqldriver.SqlDriver
}

var Driver MySQLDriver

const (
	THE_CASE string = "lower"
)

// Init function used to initialize mysql database
func Init() {
	Driver = MySQLDriver{}
	Driver.Init("mysql", configurations.Configuration.MysqlUrl)
	//utils.Info("MySQL Configs", "Url: ", configurations.Configuration.MysqlUrl)
}
