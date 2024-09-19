package jesorm

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/access6080/jesorm/structures"
)

func executeFindQuery[T any](t **T, query string, db *sql.DB) T {
	// Use db to execute sql query

	rows, err := db.Query(query)
	if err != nil {
		return **t
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return **t
	}

	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range columns {
		valuePtrs[i] = &values[i]
	}

	if rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			return **t
		}

		elem := reflect.ValueOf(*t).Elem()
		for i, col := range columns {
			field := elem.FieldByName(col)
			if field.IsValid() && field.CanSet() {
				field.Set(reflect.ValueOf(values[i]))
			}
		}
	}
	return **t
}

func constructFindQuery(tableName string, tableColumns []structures.Column, id interface{}) string {
	columns := ""
	for i, column := range tableColumns {
		if i > 0 {
			columns += ", "
		}
		columns += column.Name
	}

	query := fmt.Sprintf("SELECT %s FROM %s WHERE id = %v", columns, tableName, id)
	return query
}

func constructInsertQuery(structName func() string, cols []string) string {
	columnNames := ""
	columnValues := ""
	for i, col := range cols {
		if i > 0 {
			columnNames += ", "
			columnValues += ", "
		}
		columnNames += col
		columnValues += "?"
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", structName(), columnNames, columnValues)
	return query
}

func insertionColumns(t reflect.Type) []string {
	var columns []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		columns = append(columns, field.Name)
	}
	return columns
}

func executeInsertQuery(value any, query string, dB *sql.DB) error {
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	values := make([]interface{}, val.NumField())
	for i := 0; i < val.NumField(); i++ {
		values[i] = val.Field(i).Interface()
	}

	_, err := dB.Exec(query, values...)
	return err
}
