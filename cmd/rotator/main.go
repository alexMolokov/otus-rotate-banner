package main

import (
	"flag"
	"fmt"
	"os"

	// "github.com/alexMolokov/otus-rotate-banner/internal/app".
	configApp "github.com/alexMolokov/otus-rotate-banner/internal/config"
	"github.com/alexMolokov/otus-rotate-banner/internal/logger"
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

	fmt.Printf("%#v", cfg)

	_ = logger

	/*st, err := app.NewStorage(cfg)
	if err != nil {
		fmt.Printf("Can't create pool connect to storage: %v", err)
		os.Exit(1)
	}

	calendar := app.New(logg, st)
	defer calendar.Close()

	tcpAddr := fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port)
	httpAddr := fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)

	grpcServer := internalgrpc.NewServer(logg, calendar, tcpAddr)
	httpServer := internalhttp.NewHTTPEcoSystemServer(httpAddr, tcpAddr)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		logg.Info("Service GRPC calendar is running...")

		if err := grpcServer.Start(); err != nil {
			logg.Error("failed to start GRPC server: " + err.Error())
			cancel()
		}
	}()

	go func() {
		logg.Info("Service HTTP REST calendar is running...")

		if err := httpServer.Start(); err != nil {
			logg.Error("failed to start HTTP server: " + err.Error())
			cancel()
		}
	}()

	<-ctx.Done()

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	if err := httpServer.Stop(ctx); err != nil {
		logg.Error("failed to stop HTTP REST calendar service: " + err.Error())
	} else {
		logg.Info("Service HTTP REST calendar is stopped")
	}

	if err := grpcServer.Stop(ctx); err != nil {
		logg.Error("failed to stop GRPC calendar service: " + err.Error())
	} else {
		logg.Info("Service GRPC calendar is stopped")
	}*/
}
