package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
)

func RunMigrations(log *logrus.Logger, db *sql.DB) error {
	migrationsDir := "/app/migrations"

	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		log.Errorf("Migrations directory not found: %s", migrationsDir)
		migrationsDir = "./migrations"
	}

	log.Infof("Applying migrations from: %s", migrationsDir)

	err := goose.SetDialect("postgres")
	if err != nil {
		log.Errorf("Failed to set database dialect: %v", err)
		return fmt.Errorf("failed to set database dialect: %w", err)
	}
	log.Debug("Database dialect set to postgres")

	goose.SetTableName("goose_migrations")
	log.Debug("Migration table name set to: goose_migrations")

	log.Infof("Applying migrations from: %s", migrationsDir)
	if err := goose.Up(db, migrationsDir); err != nil {
		log.Errorf("Failed to apply migrations: %v", err)
		return fmt.Errorf("failed to apply migrations from %s: %w", migrationsDir, err)
	}

	version, err := goose.GetDBVersion(db)
	if err != nil {
		log.Warnf("Failed to get migration version: %v", err)
	} else {
		log.Infof("Current migration version: %d", version)
	}

	log.Info("Migrations completed successfully")
	return nil
}

// func findGoModRoot() string {
// 	dir, err := os.Getwd()
// 	if err != nil {
// 		return ""
// 	}

// 	for {
// 		goModPath := filepath.Join(dir, "go.mod")
// 		if _, err := os.Stat(goModPath); err == nil {
// 			return dir
// 		}

// 		parent := filepath.Dir(dir)
// 		if parent == dir {
// 			break
// 		}
// 		dir = parent
// 	}

// 	return ""
// }
