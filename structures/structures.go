package structures

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type Config struct {
	DriverName	string
	DSN			string
}

type DB struct {
	Db		 *sql.DB
    Config   Config
}

type Model struct {
    DB      *DB
    Models  ModelMap
}

type Table struct {
    Name    string
    Columns []Column
}

type Column struct {
    Name        string  `json:"name"`
    Gotype      string  `json:"goType"`
    Sqltype     string  `json:"sqlType"`
    Constraints []string  `json:"constraints"`
}

func (c *Column) IsEqual(old Column) bool {
    return c.Name == old.Name && 
           c.Gotype == old.Gotype &&
           c.Sqltype == old.Sqltype &&
           TestEq(c.Constraints, old.Constraints)

}

type Schema struct {
    SchemaName  string `json:"schemaName"`
    Columns     []Column `json:"columns"`
}

type ModelMap map[string][]Column


func (m *Model) Add(table interface{}) error {
    t := reflect.TypeOf(table)

    if t.Kind() != reflect.Struct {
        return fmt.Errorf("table must be a struct")
    }

    tableName := t.Name()
    columns := CreateColumn(table)


    m.Models[tableName] = columns

    return nil
}


func CreateColumn(table interface{}) []Column {
	var cs []Column
	t := reflect.TypeOf(table)

	for i := 0; i < t.NumField(); i++ {
		var c Column
        field := t.Field(i)
		c.Name = field.Name
		c.Gotype = field.Type.Name()
		if typeTag, ok := field.Tag.Lookup("type"); ok {
			c.Sqltype = typeTag
		}

		if consTag, ok := field.Tag.Lookup("constraints"); ok {
			c.Constraints = strings.Split(consTag, ",")
		}

		cs = append(cs, c)
    }

	return cs
}

func TestEq(a, b []string) bool {
    if len(a) != len(b) {
        return false
    }
    for i, v := range a {
        if v != b[i] {
            return false
        }
    }
    return true
}

