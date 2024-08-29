package jesorm

import (
	"database/sql"
	"fmt"
	"reflect"
)

type Config struct {
	DriverName	string
	DSN			string
}

type DB struct {
	db		 *sql.DB
    config   Config
}

type Model struct {
    DB      *DB
    models  model
}

type column struct {
    Name        string  `json:"name"`
    Gotype      string  `json:"goType"`
    Sqltype     string  `json:"sqlType"`
    Constraints []string  `json:"constraints"`
}

type schema struct {
    SchemaName  string `json:"schemaName"`
    Columns     []column `json:"columns"`
}

type model map[string][]column

func (m *Model) Add(table interface{}) error {
    t := reflect.TypeOf(table)

    if t.Kind() != reflect.Struct {
        return fmt.Errorf("table must be a struct")
    }

    tableName := t.Name()
    columns := createColumn(table)


    m.models[tableName] = columns

    return nil
}

