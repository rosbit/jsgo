/**
 * DB client module
 * Rosbit Xu <me@rosbit.cn>
 * Dec. 13, 2018
 */
package mod_db

import (
	js "github.com/rosbit/duktape-bridge/duk-bridge-go"
	"fmt"
)

const (
	MYSQL   = "mysql"
	SQLITE3 = "sqlite3"
)

type DBModule struct {
	env *js.JSEnv
}

func NewDBModule(ctx *js.JSEnv) interface{} {
	return &DBModule{ctx}
}

func getItem(options map[string]interface{}, item string, required bool) (interface{}, error) {
	i, ok := options[item]
	if !ok {
		if required {
			return nil, fmt.Errorf("item \"%s\" not found in options", item)
		}
		return nil, nil
	}
	return i, nil
}

func getString(options map[string]interface{}, item string, required bool) (string, error) {
	s, err := getItem(options, item, required)
	if err != nil {
		return "", err
	}
	if s == nil {
		return "", nil
	}
	ss, ok := s.(string)
	if !ok {
		return "", fmt.Errorf("\"%s\" must be string", item)
	}
	return ss, nil
}

func getInt(options map[string]interface{}, item string, required bool) (int, error) {
	i, err := getItem(options, item, required)
	if err != nil {
		return 0, err
	}
	if i == nil {
		return 0, nil
	}
	ii, ok := i.(float64)
	if !ok {
		return 0, fmt.Errorf("\"%s\" must be integer", item)
	}
	return int(ii), nil
}

func generateDSN(options map[string]interface{}) (string, string, error) {
	dbType, _ := getString(options, "type", false)
	switch dbType {
	case MYSQL, "":
		host, err := getString(options, "host", true)
		if err != nil {
			return "", "", err
		}

		port, err := getInt(options, "port", false)
		if err != nil {
			return "", "", err
		}
		switch {
		case port < 0:
			return "", "", fmt.Errorf("\"port\" must be greater than 0")
		case port == 0:
			port = 3306
		}

		user, err := getString(options, "user", true)
		if err != nil {
			return "", "", err
		}

		password, err := getString(options, "password", false)
		if err != nil {
			return "", "", err
		}

		dbName, err := getString(options, "database", false)
		if err != nil {
			return "", "", err
		}

		cs, err := getString(options, "charset", false)
		if err != nil {
			return "", "", err
		}
		if cs == "" {
			cs = "utf8"
		}

		var dataSource string
		if host[0] == '/' {
			dataSource = fmt.Sprintf("%s:%s@unix(%s)/%s?%s", user, password, host, dbName, cs)
		} else {
			dataSource = fmt.Sprintf("%s:%s@(%s:%d)/%s?%s", user, password, host, port, dbName, cs)
		}
		return MYSQL, dataSource, nil
	case SQLITE3:
		dbName, err := getString(options, "db", true)
		if err != nil {
			return "", "", err
		}
		return SQLITE3, dbName, nil
	default:
		return "", "", fmt.Errorf("Unknow database type: %s", dbType)
	}
}

func (m *DBModule) CreateConnection(options map[string]interface{}) (*dbConn, error) {
	if options == nil || len(options) == 0 {
		return nil, fmt.Errorf("options required to connect to create connection")
	}
	dbType, dataSource, err := generateDSN(options)
	if err != nil {
		return nil, err
	}
	return createConnctionModule(m.env, dbType, dataSource)
}

