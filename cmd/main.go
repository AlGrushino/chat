package main

import (
	"chat/internal/handlers"
	"chat/internal/repository"
	"chat/internal/service"
	"chat/pkg/db"
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	defer func() {
		log.Info("Closing database connection...")
		if err := db.Close(gormDB); err != nil {
			log.Errorf("Error closing database: %v", err)
		}
	}()

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatal("Failed to get sql.DB:", err)
	}
	defer func() {
		log.Info("Closing SQL connection...")
		if err := sqlDB.Close(); err != nil {
			log.Errorf("Error closing SQL connection: %v", err)
		}
	}()

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

	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler.GetMux(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		log.WithFields(logrus.Fields{
			"address":       server.Addr,
			"read_timeout":  server.ReadTimeout,
			"write_timeout": server.WriteTimeout,
			"idle_timeout":  server.IdleTimeout,
		}).Info("HTTP server starting")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	sig := <-stop
	log.WithField("signal", sig.String()).Info("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Info("Shutting down HTTP server gracefully...")

	if err := server.Shutdown(ctx); err != nil {
		log.Errorf("HTTP server shutdown error: %v", err)

		log.Warn("Forcing server closure...")
		if err := server.Close(); err != nil {
			log.Errorf("Failed to force close server: %v", err)
		}
	}

	log.Info("HTTP server stopped successfully")
	log.Info("Application shutdown complete")
}
