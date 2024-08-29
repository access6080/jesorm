package jesorm

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

func generateCreateQuery(tableName string, cols []column) string {
	var queryBuilder strings.Builder
    queryBuilder.WriteString(fmt.Sprintf("CREATE TABLE %s (", lower(tableName)))
    count := 0
    for _, c := range cols {		

        queryBuilder.WriteString(
			fmt.Sprintf("%s %s %s", lower(c.Name), 
				columnType(c.Gotype, c.Sqltype), 
				constraints(c.Constraints)))
        
		if count < len(cols) - 1 {
            queryBuilder.WriteString(", ")
        }
        count++
    }

    queryBuilder.WriteString(");")
    return queryBuilder.String()
}

func sqlize(goType string) string {
	switch {
	case goType == "int":
		return "INTEGER"
	case goType == "string":
		return "TEXT"
	case goType == "bool":
		return "BOOLEAN"
	}
	
	return ""
}

func lower(str string) string {
	return strings.ToLower(str)
}

func columnType(goType string, sqlType string) string {
	if sqlType == "" {
		return sqlize(goType)
	}

	return sqlType
}

func constraints(cns []string) string {
	var constraintBuilder strings.Builder

	for _, constraint := range cns {
		switch constraint {
		case "PrimaryKey":
			constraintBuilder.WriteString(" PRIMARY KEY ")
		case "NotNull":
			constraintBuilder.WriteString(" NOT NOLL ")
		case "Unique":
			constraintBuilder.WriteString(" UNIQUE ")
		default:
			
		}
	}
	return constraintBuilder.String()
}

// func formatForeignKey(column string, reference string) string {
//     return fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s", column, reference)
// }

func createColumn(table interface{}) []column {
	var cs []column
	t := reflect.TypeOf(table)

	for i := 0; i < t.NumField(); i++ {
		var c column
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

func createTable(db DB, tableName string, cs []column) error {

	query := generateCreateQuery(tableName, cs)
	db.db.Exec(query)

	return nil
}


func createOrmBaseDirectory(basepath string) error {
	return os.MkdirAll(basepath, os.ModePerm)
}

func getLastMigration(basepath string, currentMigration time.Time) (string, error) {
	// Check the migration directory for the last folder created before currentMigration
	migrationsPath := filepath.Join(basepath)
	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		return "", err
	}

	var lastMigrationTime time.Time
	var lastMigrationFolder string

	for _, file := range files {
		if file.IsDir() {
			// Parse the folder name to extract the date-time information
			folderTime, err := time.Parse("060102-150405", file.Name())
			if err != nil {
				continue // Skip folders that don't match the date-time pattern
			}

			// Compare the extracted date-time with currentMigration
			if folderTime.Before(currentMigration) && folderTime.After(lastMigrationTime) {
				lastMigrationTime = folderTime
				lastMigrationFolder = file.Name()
			}
		}
	}

	return lastMigrationFolder, nil
}

func generateSchema(m model, migrationPath string) error {
	for tableName, cols := range m {
		file, err := createSchemafile(tableName, migrationPath)
		if err != nil {
			return err
		}

		func() {
			defer func() {
                if err := file.Close(); err != nil {
                    log.Printf("Error closing file for table %s: %v", tableName, err)
                }
            }()

			s := schema{
				SchemaName: tableName,
				Columns: cols,
			}

			jsonData, err := json.Marshal(s)
			if err != nil {
				log.Printf("Error marshaling schema for table %s: %v", tableName, err)
				return 
			}

			if _, err = file.Write(jsonData); err != nil {
				log.Printf("Error writing schema to file for table %s: %v", tableName, err)
				return 
			}
		}()
	}

	return nil
}

func createSchemafile(tableName string, migrationPath string) (*os.File, error) {
	fileName := fmt.Sprintf("%s.schema.json", tableName)
	path := filepath.Join(migrationPath, fileName)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, os.ModePerm) 

	if err != nil {
		return nil, err
	}

	return file, nil
}
