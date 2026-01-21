package db

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Host     string
	User     string
	Password string
	DBname   string
	Port     string
	SSLMode  string
	Timezone string
}

func GetConfig(log *logrus.Logger) *Config {
	log.Info("Getting config from env")

	cfg := Config{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		DBname:   os.Getenv("DB_NAME"),
		Port:     os.Getenv("DB_PORT"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		Timezone: os.Getenv("DB_TIMEZONE"),
	}
	return &cfg
}

func GormInit(log *logrus.Logger, cfg *Config) (*gorm.DB, error) {
	log.Info("Initializing GORM database connection")

	dsn := getDSN(log, cfg)

	gormLogger := logger.New(
		log,
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 gormLogger,
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})

	if err != nil {
		log.Errorf("Failed to connect to database: %v", err)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Errorf("Failed to get sql.DB: %v", err)
		return nil, fmt.Errorf("Failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(30 * time.Minute)

	log.Info("Database connection established successfully")
	return db, nil
}

func getDSN(log *logrus.Logger, cfg *Config) string {
	log.Info("Getting DSN")
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.DBname,
		cfg.Port,
		cfg.SSLMode,
		cfg.Timezone,
	)
}

// func HealthCheck(db *gorm.DB) error {
// 	sqlDB, err := db.DB()
// 	if err != nil {
// 		return err
// 	}
// 	return sqlDB.Ping()
// }

func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
