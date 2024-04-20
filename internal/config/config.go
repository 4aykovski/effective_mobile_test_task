package config

import (
	"fmt"
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Postgres    PostgresConfig
	HTTP        HTTPConfig
	CarsInfoApi CarsInfoApiConfig
	Env         string `env:"ENV"`
}

type PostgresConfig struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     int    `env:"POSTGRES_PORT"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASS"`
	Database string `env:"POSTGRES_DB"`
	DSN      string
}

type HTTPConfig struct {
	Host        string `env:"HTTP_HOST"`
	Port        string `env:"HTTP_PORT"`
	Address     string
	Timeout     time.Duration `env:"HTTP_TIMEOUT"`
	IdleTimeout time.Duration `env:"HTTP_IDLE_TIMEOUT"`
}

type CarsInfoApiConfig struct {
	Host     string `env:"CARS_INFO_API_HOST"`
	BasePath string `env:"CARS_INFO_API_BASE_PATH"`
	Scheme   string `env:"CARS_INFO_API_SCHEME"`
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

	cfg.Postgres.DSN = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Database)

	cfg.HTTP.Address = fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)

	return &cfg
}
