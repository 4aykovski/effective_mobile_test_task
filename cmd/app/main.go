package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/4aykovski/effective_mobile_test_task/internal/config"
	v1 "github.com/4aykovski/effective_mobile_test_task/internal/net/v1"
	"github.com/4aykovski/effective_mobile_test_task/internal/repository/postgres"
	"github.com/4aykovski/effective_mobile_test_task/internal/service/carinfoservice"
	"github.com/4aykovski/effective_mobile_test_task/internal/service/carservice"
	"github.com/4aykovski/effective_mobile_test_task/internal/service/ownerservice"
	"github.com/4aykovski/effective_mobile_test_task/migrations"
	"github.com/4aykovski/effective_mobile_test_task/pkg/client"
	"github.com/4aykovski/effective_mobile_test_task/pkg/client/carinfo"
	postgresdb "github.com/4aykovski/effective_mobile_test_task/pkg/database/postgres"
	"github.com/4aykovski/effective_mobile_test_task/pkg/logger"
)

// @title Effective Mobile Test Task - Cars Catalog
// @version 1.0

// @host localhost:8080
// @BasePath /api/v1

// @schemes http https

func main() {
	cfg := config.MustLoad()

	log, err := logger.New(cfg.Env)
	if err != nil {
		log.Error("Failed to initialize logger", slog.String("error", err.Error()))
		os.Exit(1)
	}
	log.Debug("Logger initialized", slog.String("env", cfg.Env))

	postgresDB, err := postgresdb.New(cfg.Postgres.DSN)
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

	carRepo := postgres.NewCarRepository(postgresDB)
	ownerRepo := postgres.NewOwnerRepository(postgresDB)
	log.Debug("Repositories initialized")

	httpClient := client.NewHTTPClient(cfg.CarsInfoApi.Host, cfg.CarsInfoApi.BasePath, cfg.CarsInfoApi.Scheme, http.Client{})
	carInfoClient := carinfo.NewClient(httpClient)
	log.Debug("CarInfoClient initialized")

	carService := carservice.NewCarService(carRepo)
	ownerService := ownerservice.New(ownerRepo)
	carInfoService := carinfoservice.New(carInfoClient)
	log.Debug("Services initialized")

	mux := v1.NewMux(log, carService, ownerService, carInfoService)
	log.Debug("Mux initialized")

	httpServer := &http.Server{
		Addr:    cfg.HTTP.Address,
		Handler: mux,
	}

	log.Info("Server started", slog.String("address", cfg.HTTP.Address))

	if err := httpServer.ListenAndServe(); err != nil {
		log.Error("Failed to start server", slog.String("error", err.Error()))
		os.Exit(1)
	}

	log.Error("server stopped")
}
