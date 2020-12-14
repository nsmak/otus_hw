package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/server/http"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/config.json", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := NewConfig(configFile)
	if err != nil {
		log.Fatalf("can't get config: %v", err)
	}

	logg, err := logger.New(config.Logger.Level, config.Logger.FilePath)
	if err != nil {
		log.Fatalf("can't start logger %v\n", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// storage := memorystorage.New() // TEMP: - пока не используется
	// calendar := app.New(logg, storage)  // TEMP: - пока не используется

	server := internalhttp.NewServer(internalhttp.NewTempPublic(logg), config.HTTPServer.Host, config.HTTPServer.Port, logg)

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals)

		<-signals
		signal.Stop(signals)
		cancel()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		logg.Info("stopping server...")
		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	logg.Info("starting server at " + server.Address)
	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		os.Exit(1) // nolint: gocritic
	}
}
