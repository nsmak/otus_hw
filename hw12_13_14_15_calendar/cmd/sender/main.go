package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/cmd/config"
	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/mq/rabbit"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/calendar.json", "Path to configuration file")
}

func main() {
	flag.Parse()

	cfg, err := config.NewSender(configFile)
	if err != nil {
		log.Fatalf("can't get config: %v", err)
	}

	logg, err := logger.New(cfg.Logger.Level, cfg.Logger.FilePath)
	if err != nil {
		log.Fatalf("can't start logger %v\n", err)
	}

	rmq, err := rabbit.NewConsumer(cfg.RabbitMQ)
	if err != nil {
		log.Fatalf("can't create consumer: %v", err)
	}

	err = rmq.OpenChannel()
	if err != nil {
		log.Fatalf("can't open channel: %v", err)
	}

	err = rmq.BeginConsume()
	if err != nil {
		log.Fatalf("can't begin consume: %v", err)
	}

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals)

		<-signals
		signal.Stop(signals)

		err := rmq.CloseChannel()
		if err != nil {
			logg.Error("can't close channel", logg.String("msg", err.Error()))
		}

		err = rmq.CloseConn()
		if err != nil {
			logg.Error("can't close connection", logg.String("msg", err.Error()))
		}
	}()

	for msg := range rmq.Get() {
		if msg.Err != nil {
			log.Printf("got message with error: %s\n", msg.Err.Error())
			continue
		}
		fakeSendNotification(logg, msg.Notif)
	}
}

func fakeSendNotification(logg *logger.Logger, notif app.MQEventNotification) {
	logg.Info(
		"got message",
		logg.String("event_id", notif.EventID),
		logg.String("title", notif.Title),
		logg.Int64("date", notif.Date),
		logg.String("user_id", notif.UserID),
	)
}
