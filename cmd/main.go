package main

import (
	"chat/internal/repository"
	"chat/pkg/db"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	logFile, err := os.OpenFile("../logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	log.Info("Loading .env")
	err = godotenv.Load("../.env")
	if err != nil {
		log.Fatal(err)
	}

	cfg := db.GetConfig(log)
	gormDB, err := db.GormInit(log, cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close(gormDB)

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatal("Failed to get sql.DB:", err)
	}
	defer sqlDB.Close()

	log.Info("Running database migrations...")
	if err := db.RunMigrations(log, sqlDB); err != nil {
		log.Fatal("Migrations failed:", err)
	}
	log.Info("Migrations completed successfully")

	repo := repository.NewRepository(gormDB)
	log.Info("Creating repository")
	if err != nil {
		log.Fatalf("Failed to create repository:", err)
	}

	log.Info("Application started successfully")
}
