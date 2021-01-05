package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/cmd/config"
	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/server/grpcsrv"
	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/server/rest"
	memorystorage "github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/calendar.json", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	cfg, err := config.NewCalendar(configFile)
	if err != nil {
		log.Fatalf("can't get config: %v", err)
	}

	logg, err := logger.New(cfg.Logger.Level, cfg.Logger.FilePath)
	if err != nil {
		log.Fatalf("can't start logger %v\n", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	calendar := app.New(logg, startStorageService(ctx, cfg.Database))
	restServer := rest.NewServer(rest.NewAPI(calendar), cfg.RestServer.Host, cfg.RestServer.Port, logg)
	grpcServer := grpcsrv.NewServer(grpcsrv.NewAPI(calendar), cfg.GrpcServer.Host, cfg.GrpcServer.Port, logg)

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals)

		<-signals
		signal.Stop(signals)
		cancel()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		logg.Info("stopping rest server...")
		if err := restServer.Stop(ctx); err != nil {
			logg.Error("failed to stop rest server: " + err.Error())
		}

		logg.Info("stopping gRPC server...")
		if err := grpcServer.Stop(); err != nil {
			logg.Error("failed to stop gRPC server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		startRESTServer(ctx, restServer, logg)
	}()

	go func() {
		defer wg.Done()
		startGRPCServer(ctx, grpcServer, logg)
	}()
	wg.Wait()
}

func startRESTServer(ctx context.Context, s *rest.Server, logg app.Logger) {
	logg.Info("starting REST server at " + s.Address)
	if err := s.Start(ctx); err != nil {
		log.Fatalf("failed to start rest server: " + err.Error())
	}
}

func startGRPCServer(ctx context.Context, s *grpcsrv.Server, logg app.Logger) {
	logg.Info("starting gRPC server at " + s.Address)
	if err := s.Start(ctx); err != nil {
		log.Fatalf("failed to start gRPC server: " + err.Error())
	}
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
