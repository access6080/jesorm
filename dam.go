package jesorm

import (
	"fmt"
	"reflect"
)



func FindByID[T any](model Model, tableStruct T, id interface{}) (T, error) {
	// Get able Name from struct type
	table := new(T)
	tableType := reflect.TypeOf(table)

	if tableType.Kind() != reflect.Struct {
		return *table, fmt.Errorf("tableStruct must be a struct")
	}

	tableName := tableType.Name()
	tableColumns := model.Models[tableName]

	query := constructFindQuery(tableName, tableColumns, id)

	return executeFindQuery(&table, query, model.DB.Db), nil
}

func FindAll[T any](model Model, tableStruct T, condition Condition) error {
	panic("not implemented")
} 

func InsertOne(model Model, value any) error {
	// The the name of the value struct
	structName := reflect.TypeOf(value).Name
	cols := insertionColumns(reflect.TypeOf(value))

	query := constructInsertQuery(structName, cols)

	return executeInsertQuery(value, query, model.DB.Db)
}

func InsertMany(model, values []any) error {
	panic("not implemented")
}


func DeleteByID[T any](model Model, tableStruct T, id interface{}) error {
	panic("not implemented")
}

func Delete[T any](model Model, tableStruct T, condition Condition) error {
	panic("not implemented")
}

func Update[T any](model Model, tableStruct T, condition Condition) error {
	panic("not implemented")
}

func UpdateByID[T any](model Model, tableStruct T, id interface{}) error {
	panic("not implemented")
}