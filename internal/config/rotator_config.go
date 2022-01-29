package config

import (
	"context"
	"time"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
)

type RotatorConfig struct {
	Logger   LoggerConf
	DriverDB string `config:"driverDb"`
	DB       DBConf
	GRPC     GRPCConf
}

func NewRotatorConfig(fileName string) (*RotatorConfig, error) {
	loader := confita.NewLoader(
		file.NewBackend(fileName),
		env.NewBackend(),
	)

	cfg := DefaultRotatorConfig()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := loader.Load(ctx, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func DefaultRotatorConfig() *RotatorConfig {
	return &RotatorConfig{
		Logger: LoggerConf{
			Level:    "DEBUG",
			Encoding: "json",
			Output:   "log_rotator.json",
		},
		DB: DBConf{
			Driver:            "postgres",
			Host:              "localhost",
			Port:              5432,
			Name:              "banner",
			User:              "user",
			Password:          "pass",
			MaxConnectionPool: 5,
			SslMode:           "disable",
		},
		DriverDB: "postgres",
		GRPC: GRPCConf{
			Host: "localhost",
			Port: 9013,
		},
	}
}
