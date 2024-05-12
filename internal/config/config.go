package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env string `env:"ENV" env-required:"true"`
	HTTPServer
	DB PostgresConfig
}

type PostgresConfig struct {
	User     string `env:"PG_USER" env-required:"true"`
	Password string `env:"PG_PASSWORD" env-required:"true"`
	Host     string `env:"PG_HOST" env-required:"true"`
	Port     int    `env:"PG_EXTERNAL_PORT" env-required:"true"`
	UseSSL   string   `env:"PG_USE_SSL" env-required:"true"`
	Name     string `env:"PG_DB_NAME" env-required:"true"`
}

type HTTPServer struct {
	Host string `env:"HOST" env-required:"true"`
	Port int    `env:"PORT" env-required:"true"`
}

func Load() (*Config, error) {
	const op = "config.Load"

	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &cfg, nil
}
