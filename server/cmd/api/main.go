package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ReynoldArun09/school-procurement-management/internal/config"
	"github.com/ReynoldArun09/school-procurement-management/internal/database"
	"github.com/ReynoldArun09/school-procurement-management/internal/httpserver"
	"github.com/ReynoldArun09/school-procurement-management/internal/logger"
	"github.com/ReynoldArun09/school-procurement-management/internal/models"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	// 1. Initialize environment variables.
	cfg, err := config.Load()

	if err != nil {
		log.Printf("Configuration error: %v", err)
		os.Exit(1)
	}

	// 2. Initialize structured logging.
	appLogger := logger.New(cfg.AppEnv)

	// 3. Create a context that is cancelled by termination signals.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 4. Initialize mysql.

	db, err := database.NewMySQL(ctx, cfg.Mysql, appLogger)

	if err := db.AutoMigrate(&models.Admin{}, &models.AcademicYear{}); err != nil {
		appLogger.Error("auto-migrate database models", "error", err)
	}

	if err != nil {
		appLogger.Error("database intialization failed", "error", err)
		os.Exit(1)
	}

	sqlDB, err := db.DB()
	if err != nil {
		appLogger.Error("failed to get SQL connection pool", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			appLogger.Error("close mysql connection pool", "error", err)
		}
	}()

	// 5. Initialize the HTTP router.
	mux := http.NewServeMux()

	// 6. Initialize the HTTP server.
	server := httpserver.New(httpserver.Config{
		Host:         "",
		Port:         cfg.Port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}, mux, appLogger)

	// 7. Start the server.
	serverError := make(chan error, 1)

	go func() {
		appLogger.Info("HTTP server starting", "address", server.Addr)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverError <- err
		}
		serverError <- nil
	}()

	// 8. Wait for a termination signal or server failure.

	select {
	case <-ctx.Done():
		appLogger.Info("shutdown signal received")

	case err := <-serverError:
		if err != nil {
			log.Printf("HTTP server failed: %v", err)
		}
		log.Printf("HTTP server failed: %v", err)
	}

	// 9. Gracefully shut down the HTTP server.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		appLogger.Error("graceful shutdown failed", "error", err)

		if closeErr := server.Close(); closeErr != nil {
			log.Printf(
				"graceful shutdown failed: %v; force close failed: %v",
				err,
				closeErr,
			)
		}

		log.Printf("graceful shutdown: %v", err)
	}

	appLogger.Info("HTTP server stopped")

}
