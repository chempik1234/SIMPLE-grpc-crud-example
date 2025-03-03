package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	// Postgres postgres.Config `yaml:"POSTGRES" env:"POSTGRES"`
	GRPCPort int `yaml:"GRPC_PORT" env:"GRPC_PORT" env-default:"50051"`
	HTTPPort int `yaml:"HTTP_PORT" env:"HTTP_PORT" env-default:"8081"`
}

func New() (Config, error) {
	var cfg Config
	err := godotenv.Load()
	if err != nil {
		return Config{},
			fmt.Errorf("failed to read .env: %w", err)
	}
	if err = cleanenv.ReadEnv(&cfg); err != nil {
		return Config{},
			fmt.Errorf("failed to read env variables after accessing .env: %w", err)
	}
	return cfg, nil
}
