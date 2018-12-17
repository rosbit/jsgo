package mod_db

import (
	js "github.com/rosbit/duktape-bridge/duk-bridge-go"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type dbConn struct {
	jsEnv *js.JSEnv
	conn *sql.DB
	dbType string
}

func createConnctionModule(ctx *js.JSEnv, dbType string, dataSource string) (*dbConn, error) {
	conn, err := sql.Open(dbType, dataSource)
	if err != nil {
		return nil, err
	}
	conn.SetMaxOpenConns(1)
	if err = conn.Ping(); err != nil {
		defer conn.Close()
		return nil, err
	}

	dbConn := &dbConn{ctx, conn, dbType}
	return dbConn, nil
}

func (db *dbConn) Close() {
	db.conn.Close()
}

func (db *dbConn) End() {
	db.Close()
}

func (db *dbConn) Query(sql string, jsCallback *js.EcmaObject) error {
	if jsCallback == nil {
		return fmt.Errorf("callback function required")
	}
	defer db.jsEnv.DestroyEcmascriptFunc(jsCallback)

	if sql == "" {
		db.jsEnv.CallEcmascriptFunc(jsCallback, fmt.Errorf("sql required"), nil)
		return nil
	}

	rows, err := db.conn.Query(sql)
	if err != nil {
		db.jsEnv.CallEcmascriptFunc(jsCallback, err, nil)
		return nil
	}
	defer rows.Close()

	res, err := db.createResultModule(rows, db.dbType)
	if err != nil {
		db.jsEnv.CallEcmascriptFunc(jsCallback, err, nil)
		return nil
	}

	db.jsEnv.CallEcmascriptFunc(jsCallback, nil, res)
	return nil
}

func (db *dbConn) Update(updateSql string, params []interface{}, jsCallback *js.EcmaObject) error {
	if jsCallback == nil {
		return fmt.Errorf("callback function required")
	}
	defer db.jsEnv.DestroyEcmascriptFunc(jsCallback)

	if updateSql == "" {
		db.jsEnv.CallEcmascriptFunc(jsCallback, fmt.Errorf("sql required"), nil)
		return nil
	}

	var res sql.Result
	var err error
	if params == nil {
		res, err = db.conn.Exec(updateSql)
	} else {
		res, err = db.conn.Exec(updateSql, params...)
	}
	if err != nil {
		db.jsEnv.CallEcmascriptFunc(jsCallback, fmt.Errorf("sql required"), nil)
		return nil
	}

	db.jsEnv.CallEcmascriptFunc(jsCallback, nil, res)
	return nil
}

func (db *dbConn) createResultModule(rows *sql.Rows, dbType string) (*resultSet, error) {
	rs := &resultSet{rows:rows, dbType:dbType}
	rs.columns, _ = rows.Columns()
	rs.colTypes, _ = rows.ColumnTypes()
	colNum := len(rs.columns)
	rs.scanArgs = make([]interface{}, colNum)
	rs.row = make([]interface{}, colNum)
	for i := range rs.row {
		rs.scanArgs[i] = &rs.row[i]
	}
	rs.mapRow = make(map[string]interface{}, colNum)
	rs.arrRow = make([]interface{}, colNum)
	return rs, nil
}

type resultSet struct {
	rows *sql.Rows
	dbType string
	columns []string
	colTypes []*sql.ColumnType
	row []interface{}
	scanArgs []interface{}
	mapRow map[string]interface{}
	arrRow []interface{}
}

func (rs *resultSet) Fields() []string {
	return rs.columns
}

var (
	_trueVals = map[string]bool {
		"1":true,
		"y":true,
		"Y":true,
		"yes":true,
		"Yes":true,
		"YES":true,
	}
)

func recognizeMySQLType(v []byte, colType *sql.ColumnType) interface{} {
	var res interface{}
	scanType := colType.ScanType()

	switch scanType.Kind() {
	case reflect.Int8,reflect.Uint8,reflect.Int16,reflect.Uint16,reflect.Int32,reflect.Int:
		res, _ = strconv.Atoi(string(v))
	case reflect.Uint32,reflect.Uint,reflect.Int64:
		res, _ = strconv.ParseInt(string(v), 10, 64)
	case reflect.Uint64:
		res, _ = strconv.ParseUint(string(v), 10, 64)
	case reflect.Float32,reflect.Float64:
		res, _ = strconv.ParseFloat(string(v), 64)
	case reflect.Bool:
		_, res = _trueVals[string(v)]
	case reflect.Struct:
		switch scanType.Name() {
		case "NullBool":
			_, res = _trueVals[string(v)]
		case "NullFloat64":
			res, _ = strconv.ParseFloat(string(v), 64)
		case "NullInt64":
			res, _ = strconv.ParseInt(string(v), 10, 64)
		default:
			res = string(v)
		}
	case reflect.Slice:
		if colType.DatabaseTypeName() != "DECIMAL" {
			res = string(v)
		} else {
			res, _ = strconv.ParseFloat(string(v), 64)
		}
	default:
		res = string(v)
	}
	return res
}

func recognizeSQLite3Type(col interface{}) interface{} {
	switch col.(type) {
	case []byte:
		return string(col.([]byte))
	case time.Time:
		return col.(time.Time).Format("2006-01-02 15:04:05.000")
	default:
		return col
	}
}

func (rs *resultSet) Next() map[string]interface{} {
	if rs.rows.Next() {
		rs.rows.Scan(rs.scanArgs...)
		switch rs.dbType {
		case MYSQL:
			var v []byte
			for i, col := range rs.row {
				colName := rs.columns[i]
				if col != nil {
					v = col.([]byte)
					rs.mapRow[colName] = recognizeMySQLType(v, rs.colTypes[i])
				} else {
					rs.mapRow[colName] = nil
				}
			}
		case SQLITE3:
			for i, col := range rs.row {
				colName := rs.columns[i]
				if col != nil {
					// v := reflect.ValueOf(col)
					// fmt.Printf("%s value type: %v, type: %v\n", colName, v.Kind(), v.Type())
					rs.mapRow[colName] = recognizeSQLite3Type(col)
					switch col.(type) {
					case []byte:
						rs.mapRow[colName] = string(col.([]byte))
					case time.Time:
						rs.mapRow[colName] = col.(time.Time).Format("2006-01-02 15:04:05.000")
					default:
						rs.mapRow[colName] = col
					}
				} else {
					rs.mapRow[colName] = nil
				}
			}
		}
		return rs.mapRow
	}
	return nil
}

func (rs *resultSet) NextRow() []interface{} {
	if rs.rows.Next() {
		rs.rows.Scan(rs.scanArgs...)
		switch rs.dbType {
		case MYSQL:
			var v []byte
			for i, col := range rs.row {
				if col != nil {
					v = col.([]byte)
					rs.arrRow[i] = recognizeMySQLType(v, rs.colTypes[i])
				} else {
					rs.arrRow[i] = nil
				}
			}
		case SQLITE3:
			for i, col := range rs.row {
				if col != nil {
					rs.arrRow[i] = recognizeSQLite3Type(col)
				} else {
					rs.arrRow[i] = nil
				}
			}
		}
		return rs.arrRow
	}
	return nil
}
