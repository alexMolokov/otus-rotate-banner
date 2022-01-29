package config

import (
	"context"
	"time"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
)

type SchedulerConfig struct {
	DriverDB string
	Logger   LoggerConf
	DB       DBConf
	Queue    QueueConf
}

func NewSchedulerConfig(fileName string) (*SchedulerConfig, error) {
	loader := confita.NewLoader(
		file.NewBackend(fileName),
		env.NewBackend(),
	)
	cfg := DefaultSchedulerConfig()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := loader.Load(ctx, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func DefaultSchedulerConfig() *SchedulerConfig {
	return &SchedulerConfig{
		Logger: LoggerConf{
			Level:    "INFO",
			Encoding: "json",
			Output:   "log_scheduler.json",
		},
		DB: DBConf{
			Driver:            "postgres",
			Host:              "localhost",
			Port:              5432,
			Name:              "calendar",
			User:              "user",
			Password:          "password",
			MaxConnectionPool: 2,
			SslMode:           "disable",
		},
		Queue: QueueConf{
			URI:            "amqp://guest:guest@localhost:5672/",
			Name:           "reminder_events",
			ExchangeName:   "reminder",
			ExchangeType:   "direct",
			ReInitDelay:    2,
			ReconnectDelay: 5,
			ResendDelay:    5,
		},
	}
}
