package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/server/grpcsrv"
	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/server/rest"
	memorystorage "github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/storage/sql"
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

	calendar := app.New(logg, startStorageService(ctx, config.Database))
	restServer := rest.NewServer(rest.NewAPI(calendar), config.RestServer.Host, config.RestServer.Port, logg)
	grpcServer := grpcsrv.NewServer(grpcsrv.NewAPI(calendar), config.GrpcServer.Host, config.GrpcServer.Port, logg)

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

func startStorageService(ctx context.Context, config DBConf) app.Storage {
	var s app.Storage
	if config.InMem {
		s = memorystorage.New()
	} else {
		sqlStore, err := sqlstorage.New(ctx, config.Username, config.Password, config.Address, config.DBName)
		if err != nil {
			log.Fatalf("failed to start storage connection: " + err.Error())
		}
		s = sqlStore
	}
	return s
}
