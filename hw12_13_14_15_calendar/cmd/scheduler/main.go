package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/cmd/config"
	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/mq/rabbit"
	memorystorage "github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/scheduler.json", "Path to configuration file")
}

func main() {
	flag.Parse()

	cfg, err := config.NewScheduler(configFile)
	if err != nil {
		log.Fatalf("can't get config: %v", err)
	}

	logg, err := logger.New(cfg.Logger.Level, cfg.Logger.FilePath)
	if err != nil {
		log.Fatalf("can't start logger %v\n", err)
	}

	producer, err := rabbit.NewProducer(cfg.RabbitMQ)
	if err != nil {
		log.Fatalf("can't create consumer: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scheduler := app.NewScheduler(
		logg,
		startStorageService(ctx, cfg.Database),
		producer,
		time.Duration(cfg.IntervalInSec)*time.Second,
	)

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)

		<-signals
		signal.Stop(signals)
		err := producer.CloseConn()
		if err != nil {
			logg.Error("can't close connection", logg.String("msg", err.Error()))
		}
		cancel()
	}()

	scheduler.Run(ctx)
}

func startStorageService(ctx context.Context, cfg config.DBConf) app.Storage {
	var s app.Storage
	if cfg.InMem {
		s = memorystorage.New()
	} else {
		sqlStore, err := sqlstorage.New(ctx, cfg.Username, cfg.Password, cfg.Address, cfg.DBName)
		if err != nil {
			log.Fatalf("failed to start storage connection: " + err.Error())
		}
		s = sqlStore
	}
	return s
}
