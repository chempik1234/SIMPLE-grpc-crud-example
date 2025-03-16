package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"yandexLyceumTheme3gRPC/pkg/postgres"
)

type Config struct {
	Postgres postgres.Config `yaml:"POSTGRES" env:"POSTGRES" env-prefix:""`
	GRPCPort int             `yaml:"GRPC_PORT" env:"GRPC_PORT" env-default:"50051"`
	HTTPPort int             `yaml:"HTTP_PORT" env:"HTTP_PORT" env-default:"8081"`
}

func New() (Config, error) {
	var cfg Config
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("warning: failed to read .env: %s", err)
	}
	if err = cleanenv.ReadEnv(&cfg); err != nil {
		return Config{},
			fmt.Errorf("failed to read env variables after accessing .env: %w", err)
	}
	return cfg, nil
}
