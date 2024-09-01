package jesorm

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/access6080/jesorm/helpers"
	"github.com/access6080/jesorm/structures"
	_ "github.com/go-sql-driver/mysql" // MySQL driver
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

func  Init(config structures.Config) (*structures.Model, error) {
	switch config.DriverName {
	case "mysql", "sqlite3", "postgres":
		db, err := sql.Open(config.DriverName, config.DSN)
		if err != nil {
			return nil,  err
		}
	
		if err := db.Ping(); err != nil {
			return nil, err
		}
		
		
		return &structures.Model{
				DB: &structures.DB{
					Db: db,
					Config: config,
				},
				Models: make(structures.ModelMap),
			}, nil
    default:
        return nil, fmt.Errorf("unsupported driver: %s", config.DriverName)
	}
}

func AutoMigrate(m structures.Model) error {
	basepath := filepath.Join("jesorm", "schemas", "migrations")
	if err := helpers.CreateOrmBaseDirectory(basepath); err != nil {
		return err
	}

	// Check migration needed
	currentMigration := time.Now()

	lastMigrationFolder, err := helpers.GetLastMigration(basepath, currentMigration)
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

		if err := helpers.GenerateSchema(m.Models, path); err != nil {
			return err
		}

		for tableName, cols := range m.Models {
			if err := helpers.CreateTable(*m.DB, tableName, cols); err != nil {
				return err
			}
		}

		return nil
	}

	// A schema already exists, get the last schema and compare it to models.
	oldModels, err := helpers.GetLastModels(lastMigrationFolder)
	if err != nil {
		return err
	}

	migrate, err := helpers.CompareModels(oldModels, m.Models)
	if err != nil {
		return err
	}

	// Perform migration if any changes exists
	if migrate {
		
	}

	return nil
}



