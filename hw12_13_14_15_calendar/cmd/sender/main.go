package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	_ "github.com/jackc/pgx/v4/stdlib" // nolint: gci

	"github.com/jmoiron/sqlx"
	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/cmd/config"
	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/mq/rabbit"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/sender.json", "Path to configuration file")
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
		signal.Notify(signals, os.Interrupt)

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
		fakeSendNotification(logg, cfg.Database, msg.Notif)
	}
}

func fakeSendNotification(logg *logger.Logger, cfg config.DBConf, notif app.MQEventNotification) { // только для теста
	logg.Info(
		"got message",
		logg.String("event_id", notif.EventID),
		logg.String("title", notif.Title),
		logg.Int64("date", notif.Date),
		logg.String("user_id", notif.UserID),
	)

	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", cfg.Username, cfg.Password, cfg.Address, cfg.DBName)
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		logg.Error("can't open sql", logg.String("msg", err.Error()))
	}

	err = db.Ping()
	if err != nil {
		logg.Error("can't ping sql", logg.String("msg", err.Error()))
	}

	_, err = db.Exec(`INSERT INTO notification (id, title, start_date) 
			VALUES ($1, $2, $3)`,
		notif.EventID,
		notif.Title,
		notif.Date,
	)
	if err != nil {
		logg.Error("can't create notification in db", logg.String("msg", err.Error()))
	}
}
