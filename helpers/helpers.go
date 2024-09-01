package helpers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/access6080/jesorm/structures"
)




func generateCreateQuery(tableName string, cols []structures.Column) string {
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


func CreateTable(db structures.DB, tableName string, cs []structures.Column) error {

	query := generateCreateQuery(tableName, cs)
	db.Db.Exec(query)

	return nil
}


func CreateOrmBaseDirectory(basepath string) error {
	return os.MkdirAll(basepath, os.ModePerm)
}

func GetLastMigration(basepath string, currentMigration time.Time) (string, error) {
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

func GenerateSchema(m structures.ModelMap, migrationPath string) error {
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

			s := structures.Schema{
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

func CompareModels(oldModels map[string][]structures.Column, models structures.ModelMap) (bool, error) {
	res := false

	for tableName, cols := range models {
		oldCol, exists := oldModels[tableName]
		if !exists {
			res = true
			continue
		}

		if len(oldCol) != len(cols) {
            res = true
            continue
        }

		// for i := range cols {
        //     if !cols[i].IsEqual(oldCol[i]) {
        //         res = false
        //         break
        //     }
        // }
		
	}

	return res, nil
}

func GetLastModels(lastMigrationFolder string) (map[string][]structures.Column, error) {
	result := make(map[string][]structures.Column)

    files, err := os.ReadDir(lastMigrationFolder)
    if err != nil {
        return nil, fmt.Errorf("error reading directory %s: %w", lastMigrationFolder, err)
    }

    for _, file := range files {
        if file.IsDir() {
            continue
        }

        filePath := filepath.Join(lastMigrationFolder, file.Name())
        fileContent, err := os.ReadFile(filePath)
        if err != nil {
            return nil, fmt.Errorf("error reading file %s: %w", filePath, err)
        }

        var col []structures.Column
        if err := json.Unmarshal(fileContent, &col); err != nil {
            return nil, fmt.Errorf("error unmarshaling JSON from file %s: %w", filePath, err)
        }

		name := strings.Split(file.Name(), ".")[0]

        result[name] = col
    }

    return result, nil
}


