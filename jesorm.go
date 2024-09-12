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

func Init(config structures.Config) (*structures.Model, error) {
	switch config.DriverName {
	case "mysql", "sqlite3", "postgres":
		db, err := sql.Open(config.DriverName, config.DSN)
		if err != nil {
			return nil, err
		}

		if err := db.Ping(); err != nil {
			return nil, err
		}

		return &structures.Model{
			DB: &structures.DB{
				Db:     db,
				Config: config,
			},
			Models: make(structures.ModelMap),
		}, nil
	default:
		return nil, fmt.Errorf("unsupported driver: %s", config.DriverName)
	}
}

func AutoMigrate(m structures.Model) error {
	basePath := filepath.Join("jesorm", "schemas", "migrations")
	if err := helpers.CreateOrmBaseDirectory(basePath); err != nil {
		return err
	}

	// Check migration needed
	currentMigration := time.Now()

	lastMigrationFolder, err := helpers.GetLastMigration(basePath, currentMigration)
	if err != nil {
		return err
	}

	// First Migration
	if lastMigrationFolder == "" {
		// Create folder for current migration
		path := filepath.Join(basePath, currentMigration.Format("060102-150405"))
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
	lastMigrationFolder = filepath.Join(basePath, lastMigrationFolder)
	oldModels, err := helpers.GetLastModels(lastMigrationFolder)
	if err != nil {
		return err
	}

	migrate, migrateModels := helpers.CompareModels(oldModels, m.Models)

	// Perform migration if any changes exists
	if migrate {
		//Create new schema files
		path := filepath.Join(basePath, currentMigration.Format("060102-150405"))
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}

		newModels := make(structures.ModelMap)
		for key := range migrateModels {
			cols := m.Models[key]
			newModels[key] = append(newModels[key], cols...)
		}

		if err := helpers.GenerateSchema(newModels, path); err != nil {
			return err
		}

		helpers.PerformMigration(migrateModels, m)
	}

	return nil
}
