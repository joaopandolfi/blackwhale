package sqldriver

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/elgs/gosqljson"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	//_ "gopkg.in/rana/ora.v4"
	//_ "github.com/mattn/go-oci8"
	"github.com/joaopandolfi/blackwhale/v2/utils"
)

type RemoteSqlDriver any

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
	}
	return
}

// Force request ignoring foreign keys
func (cc SqlDriver) ForceRequest() (err error) {
	err = cc.Execute("lower", nil, "SET FOREIGN_KEY_CHECKS=0;")
	if err != nil {
		utils.Error(fmt.Sprintf("[SQLDriver][%s]- Error on FORCE REQUEST", cc.DriverName), err)
	}
	return
}

// Execute method is used for execute a SQL
func (cc SqlDriver) Execute(theCase string, output any, sqlStatement string, sqlParams ...any) (err error) {
	cc.getDB()
	data, err := gosqljson.QueryToMaps(cc.Database, toCase(theCase), sqlStatement, sqlParams...)
	if err != nil {
		utils.Error(fmt.Sprintf("[SQLDriver][%s]- Error on execute query", cc.DriverName), err)
		return
	}

	b, mErr := json.Marshal(data)
	if mErr != nil {
		utils.Error(fmt.Sprintf("[SQLDriver][%s]- Error on marshal query", cc.DriverName), mErr)
		return mErr
	}

	return json.Unmarshal(b, &output)
}

func (cc SqlDriver) Run(output any, sqlStatement string, sqlParams ...any) (err error) {
	cc.getDB()
	output, err = cc.Database.Exec(sqlStatement, sqlParams...)
	return
}

func (cc SqlDriver) QueryRow(sqlStatement string, sqlParams ...any) (row *sql.Row) {
	cc.getDB()
	row = cc.Database.QueryRow(sqlStatement, sqlParams...)
	return
}

func (cc SqlDriver) ReadDBMS(output any) (result string, err error) {
	var a int
	cc.getDB()
	_, err = cc.Database.Exec(`BEGIN DBMS_OUTPUT.GET_LINE(:lines, :status); END;`,
		sql.Named("lines", sql.Out{Dest: &result}),
		sql.Named("status", sql.Out{Dest: &a, In: true}))
	return
}

func (cc SqlDriver) QueryContext(output any, sqlStatement string, sqlParams ...any) (err error) {
	cc.getDB()
	ctx, cancel := context.WithTimeout(context.Background(), 55*time.Second)
	defer cancel()
	output, err = cc.Database.QueryContext(ctx, sqlStatement, sqlParams...)
	return err
}

func (cc SqlDriver) ExecuteAndReturnLastId(sqlStatement string, sqlParams ...any) (id int64, err error) {
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
func (cc SqlDriver) ExecuteToArray(theCase string, sqlStatement string, sqlParams ...any) (header []string, data [][]string, err error) {
	cc.getDB()
	var rows [][]any
	header, rows, err = gosqljson.QueryToArrays(cc.Database, toCase(theCase), sqlStatement, sqlParams...)

	if err != nil {
		utils.Error(fmt.Sprintf("[SQLDriver][%s]- Error on Execute query to array", cc.DriverName), err)
		return
	}

	data = make([][]string, len(rows))
	for i, row := range rows {
		data[i] = make([]string, len(row))
		for j, cell := range row {
			data[i][j] = cellToString(cell)
		}
	}

	return
}

func (cc SqlDriver) QueryToMap(theCase string, sqlStatement string, sqlParams ...any) (data []map[string]string, err error) {
	cc.getDB()
	var rows []map[string]any
	rows, err = gosqljson.QueryToMaps(cc.Database, toCase(theCase), sqlStatement, sqlParams...)

	if err != nil {
		utils.Error(fmt.Sprintf("[SQLDriver][%s]- Error on query to map", cc.DriverName), err)
		return
	}

	data = make([]map[string]string, len(rows))
	for i, m := range rows {
		data[i] = make(map[string]string, len(m))
		for k, v := range m {
			data[i][k] = cellToString(v)
		}
	}

	return
}

func toCase(theCase string) int {
	switch strings.ToLower(theCase) {
	case "lower":
		return gosqljson.Lower
	case "upper":
		return gosqljson.Upper
	case "camel":
		return gosqljson.Camel
	default:
		return gosqljson.AsIs
	}
}

func cellToString(v any) string {
	if b, ok := v.([]byte); ok {
		return string(b)
	}
	return fmt.Sprintf("%v", v)
}

// QueryToJSON - return ditectly on byte array
func (cc SqlDriver) QueryToJSON(sqlStatement string, sqlParams ...any) ([]byte, error) {
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

	tableData := make([]map[string]any, 0)

	count := len(columns)
	values := make([]any, count)
	scanArgs := make([]any, count)
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		entry := make(map[string]any)
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
