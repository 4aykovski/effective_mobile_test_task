package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/4aykovski/effective_mobile_test_task/internal/config"
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
}
