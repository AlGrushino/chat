package main

import (
	"chat/internal/handlers"
	"chat/internal/repository"
	"chat/internal/service"
	"chat/pkg/db"
	"flag"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	migrateOnly := flag.Bool("migrate", false, "Run migrations only and exit")
	flag.Parse()

	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	logFile, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	log.Info("Starting chat application")

	if err := godotenv.Load(".env"); err != nil {
		log.Warn("No .env file found, using environment variables")
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

	if *migrateOnly {
		log.Info("Running migrations only mode")
		if err := db.RunMigrations(log, sqlDB); err != nil {
			log.Fatal("Migrations failed:", err)
		}
		log.Info("Migrations completed successfully, exiting")
		return
	}

	repo := repository.NewRepository(gormDB)
	svc := service.NewService(log, repo)
	handler := handlers.NewHandler(svc, log)

	handler.InitRoutes()

	addr := ":8080"
	log.WithField("address", addr).Info("Server starting")

	if err := http.ListenAndServe(addr, handler.GetMux()); err != nil {
		log.Fatal("Server failed:", err)
	}
}
