package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Postgres PostgresConfig
	HTTP     HTTPConfig
	Env      string `env:"ENV"`
}

type PostgresConfig struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASS"`
	Database string `env:"POSTGRES_DB"`
}

type HTTPConfig struct {
	Host        string        `env:"HTTP_HOST"`
	Port        string        `env:"HTTP_PORT"`
	Timeout     time.Duration `env:"HTTP_TIMEOUT"`
	IdleTimeout time.Duration `env:"HTTP_IDLE_TIMEOUT"`
}

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err.Error())
	}

	var cfg Config
	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err.Error())
	}

	return &cfg
}
