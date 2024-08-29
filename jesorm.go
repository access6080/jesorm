package jesorm

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

func  Init(config Config) (*Model, error) {
	switch config.DriverName {
	case "mysql", "sqlite3", "postgres":
		db, err := sql.Open(config.DriverName, config.DSN)
		if err != nil {
			return nil,  err
		}
	
		if err := db.Ping(); err != nil {
			return nil, err
		}
		
		
		return &Model{
				DB: &DB{
					db: db,
					config: config,
				},
				models: make(model),
			}, nil
    default:
        return nil, fmt.Errorf("unsupported driver: %s", config.DriverName)
	}
}

func AutoMigrate(m Model) error {
	basepath := filepath.Join("jesorm", "schemas", "migrations")
	if err := createOrmBaseDirectory(basepath); err != nil {
		return err
	}

	// Check migration needed
	currentMigration := time.Now()

	lastMigrationFolder, err := getLastMigration(basepath, currentMigration)
	if err != nil {
		return err
	}

	// First Migration
	if lastMigrationFolder == "" { 
		// Create folder for current migration
		path := filepath.Join(basepath, currentMigration.Format("060102-150405"))
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}

		if err := generateSchema(m.models, path); err != nil {
			return err
		}

		for tableName, cols := range m.models {
			if err := createTable(*m.DB, tableName, cols); err != nil {
				return err
			}
		}

		return nil
	}

	// A schema already exists, get the last schema and compare it to models.
	// Perform migration if any changes exists
	

	return nil
}



