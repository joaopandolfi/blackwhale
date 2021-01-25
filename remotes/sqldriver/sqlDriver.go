package sqldriver

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/elgs/gosqljson"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	//_ "gopkg.in/rana/ora.v4"
	//_ "github.com/mattn/go-oci8"
	"github.com/joaopandolfi/blackwhale/utils"
)

type RemoteSqlDriver interface{}

type SqlDriver struct {
	DriverName string
	Url        string
	Database   *sql.DB
}

// Init function used to initialize mysql Database
func (cc *SqlDriver) Init(driverName string, url string) {
	cc.DriverName = driverName
	cc.Url = url

	var err error
	//if(cc.Database == nil) {
	cc.Database, err = sql.Open(cc.DriverName, cc.Url)
	//}

	if err != nil {
		utils.CriticalError("[SQL]- Erro ao conectar o driver "+cc.DriverName, err)
		//panic(err)
	}
}

// Check if have connection
func (cc SqlDriver) getDB() *sql.DB {
	if cc.Database == nil {
		cc.Database, _ = sql.Open(cc.DriverName, cc.Url)
		utils.Info("[SQL]- New connection created", cc.DriverName)
	}

	return cc.Database
}

func (cc SqlDriver) Close() (err error) {
	if cc.Database != nil {
		err = cc.Database.Close()

		if err != nil {
			utils.CriticalError("[SQLDriver][Close] - Error on close connection", err)
		}
		cc.Database = nil
	}

	return
}

func (cc SqlDriver) RenewConnection() (err error) {
	cc.Close()

	if cc.Database == nil {
		cc.Database, err = sql.Open(cc.DriverName, cc.Url)
	}

	if err != nil {
		utils.CriticalError("[SQL]- Erro ao conectar o driver "+cc.DriverName, err)
		//panic(err)
	}
	return
}

// Force request ignoring foreign keys
func (cc SqlDriver) ForceRequest() (err error) {
	err = cc.Execute("lower", nil, "SET FOREIGN_KEY_CHECKS=0;")
	if err != nil {
		utils.Error(fmt.Sprintf("[SQLDriver][%s]- Error on FORCE REQUEST", cc.DriverName), err)
		//panic(err)
	}
	return
}

// Execute method is used for execute a SQL
func (cc SqlDriver) Execute(theCase string, output interface{}, sqlStatement string, sqlParams ...interface{}) (err error) {
	cc.getDB()
	data, err := gosqljson.QueryDbToMapJSON(cc.Database, theCase, sqlStatement, sqlParams...)
	if err != nil {
		utils.Error(fmt.Sprintf("[SQLDriver][%s]- Error on execute query", cc.DriverName), err)
		//panic(err)
	}

	json.Unmarshal([]byte(data), &output)
	return
}

func (cc SqlDriver) Run(output interface{}, sqlStatement string, sqlParams ...interface{}) (err error) {
	cc.getDB()
	output, err = cc.Database.Exec(sqlStatement, sqlParams...)
	return
}

func (cc SqlDriver) QueryRow(sqlStatement string, sqlParams ...interface{}) (row *sql.Row) {
	cc.getDB()
	row = cc.Database.QueryRow(sqlStatement, sqlParams...)
	return
}

func (cc SqlDriver) ReadDBMS(output interface{}) (result string, err error) {
	var a int
	cc.getDB()
	_, err = cc.Database.Exec(`BEGIN DBMS_OUTPUT.GET_LINE(:lines, :status); END;`,
		sql.Named("lines", sql.Out{Dest: &result}),
		sql.Named("status", sql.Out{Dest: &a, In: true}))
	return
}

func (cc SqlDriver) QueryContext(output interface{}, sqlStatement string, sqlParams ...interface{}) (err error) {
	cc.getDB()
	ctx, cancel := context.WithTimeout(context.Background(), 55*time.Second)
	defer cancel()
	output, err = cc.Database.QueryContext(ctx, sqlStatement, sqlParams...)
	return err
}

func (cc SqlDriver) ExecuteAndReturnLastId(sqlStatement string, sqlParams ...interface{}) (id int64, err error) {
	cc.getDB()
	res, err := cc.Database.Exec(sqlStatement, sqlParams...)

	if err == nil {
		id, err = res.LastInsertId()
	} else {
		utils.Error(fmt.Sprintf("[SQLDriver][%s]- Error on query and return last id", cc.DriverName), err)
	}

	return
}

// ExecuteToArray method is used for execute a SQL
func (cc SqlDriver) ExecuteToArray(theCase string, sqlStatement string, sqlParams ...interface{}) (header []string, data [][]string, err error) {
	cc.getDB()
	header, data, err = gosqljson.QueryDbToArray(cc.Database, theCase, sqlStatement, sqlParams...)

	if err != nil {
		utils.Error(fmt.Sprintf("[SQLDriver][%s]- Error on Execute query to array", cc.DriverName), err)
		//panic(err)
	}

	return
}

func (cc SqlDriver) QueryToMap(theCase string, sqlStatement string, sqlParams ...interface{}) (data []map[string]string, err error) {
	cc.getDB()
	data, err = gosqljson.QueryDbToMap(cc.Database, theCase, sqlStatement, sqlParams...)

	if err != nil {
		utils.Error(fmt.Sprintf("[SQLDriver][%s]- Error on query to map", cc.DriverName), err)
		//panic(err)
	}

	return
}

// QueryToJSON - return ditectly on byte array
func (cc SqlDriver) QueryToJSON(sqlStatement string, sqlParams ...interface{}) ([]byte, error) {
	cc.getDB()
	rows, err := cc.Database.Query(sqlStatement, sqlParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	tableData := make([]map[string]interface{}, 0)

	count := len(columns)
	values := make([]interface{}, count)
	scanArgs := make([]interface{}, count)
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		entry := make(map[string]interface{})
		for i, col := range columns {
			v := values[i]

			b, ok := v.([]byte)
			if ok {
				entry[col] = string(b)
			} else {
				entry[col] = v
			}
		}

		tableData = append(tableData, entry)
	}

	jsonData, err := json.Marshal(tableData)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}
