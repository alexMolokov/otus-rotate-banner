package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	approtator "github.com/alexMolokov/otus-rotate-banner/internal/app/rotator"

	// "github.com/alexMolokov/otus-rotate-banner/internal/app".
	configApp "github.com/alexMolokov/otus-rotate-banner/internal/config"
	"github.com/alexMolokov/otus-rotate-banner/internal/logger"
	internalgrpc "github.com/alexMolokov/otus-rotate-banner/internal/server/grpc"
	rs "github.com/alexMolokov/otus-rotate-banner/internal/storage/rotator"
)

var configFile string

func init() {
	flag.StringVar(
		&configFile,
		"config",
		"./configs/rotator.json",
		"Path to configuration file",
	)
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	cfg, err := configApp.NewRotatorConfig(configFile)
	if err != nil {
		fmt.Printf("Can't load config: %v", err)
		os.Exit(1)
	}

	logger, err := logger.New(&cfg.Logger)
	if err != nil {
		fmt.Printf("Can't create logger: %v", err)
		os.Exit(1)
	}

	storage := rs.NewRotatorStorage(cfg.DB)
	err = storage.Connect()
	if err != nil {
		fmt.Printf("Can't create pool connect to storage: %v", err)
		os.Exit(1)
	}
	defer func() {
		err = storage.Close()
		if err != nil {
			logger.Error("failed to close pool connection to storage: " + err.Error())
		}
	}()

	app := approtator.NewAppRotator(logger, storage)

	tcpAddr := fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port)
	grpcServer := internalgrpc.NewServer(logger, app, tcpAddr)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		logger.Info("Service GRPC bannerRotator is running...")

		if err := grpcServer.Start(); err != nil {
			logger.Error("failed to start GRPC server: " + err.Error())
			cancel()
		}
	}()

	<-ctx.Done()

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	if err := grpcServer.Stop(ctx); err != nil {
		logger.Error("failed to stop GRPC bannerRotator service: " + err.Error())
	} else {
		logger.Info("Service GRPC bannerRotator is stopped")
	}
}
