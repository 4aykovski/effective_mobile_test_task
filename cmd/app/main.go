package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/4aykovski/effective_mobile_test_task/internal/config"
	"github.com/4aykovski/effective_mobile_test_task/migrations"
	"github.com/4aykovski/effective_mobile_test_task/pkg/database/postgres"
	"github.com/4aykovski/effective_mobile_test_task/pkg/logger"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	log, err := logger.New(cfg.Env)
	if err != nil {
		log.Error("Failed to initialize logger", slog.String("error", err.Error()))
		os.Exit(1)
	}
	log.Debug("Logger initialized", slog.String("env", cfg.Env))

	postgresDB, err := postgres.New(cfg.Postgres.DSN)
	if err != nil {
		log.Error("Failed to connect to postgres", slog.String("error", err.Error()))
		os.Exit(1)
	}
	log.Debug("Postgres connected", slog.String("dsn", cfg.Postgres.DSN))

	if err = migrations.RunMigrations(postgresDB.DB); err != nil {
		log.Error("Failed to run migrations", slog.String("error", err.Error()))
		os.Exit(1)
	}
	log.Info("Migrations applied")

}
