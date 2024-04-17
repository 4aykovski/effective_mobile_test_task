package logger

import (
	"fmt"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	envLocal = "local"
	envProd  = "prod"
	envDev   = "dev"
)

func New(env string) (*slog.Logger, error) {
	var log *slog.Logger

	switch env {
	case envLocal:
		options := slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		}
		log = slog.New(slog.NewTextHandler(os.Stdout, &options))
	case envDev:
		lumber := lumberjack.Logger{
			Filename:   "logs/log.log",
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     28, // days
			Compress:   true,
		}
		options := slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		}
		log = slog.New(slog.NewJSONHandler(&lumber, &options))
	case envProd:
		lumber := lumberjack.Logger{
			Filename:   "logs/log.log",
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     28, // days
			Compress:   true,
		}
		options := slog.HandlerOptions{
			Level:     slog.LevelInfo,
			AddSource: true,
		}
		log = slog.New(slog.NewJSONHandler(&lumber, &options))
	default:
		return nil, fmt.Errorf("unknown environment: %s", env)
	}

	return log, nil
}
