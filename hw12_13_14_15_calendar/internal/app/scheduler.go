package app

import (
	"context"
	"encoding/json"
	"time"
)

type Scheduler struct {
	log      Logger
	storage  Storage
	producer MQProducer
	interval time.Duration
}

func NewScheduler(logger Logger, storage Storage, producer MQProducer, interval time.Duration) *Scheduler {
	return &Scheduler{log: logger, storage: storage, producer: producer, interval: interval}
}

func (s *Scheduler) Run(ctx context.Context) {
	doneCh := make(chan struct{})
	go startWorker(ctx, doneCh, s.interval, func() {
		s.publishNotificationMessage(ctx)
	})
	go startWorker(ctx, doneCh, s.interval, func() {
		s.clearEvents(ctx)
	})
	<-doneCh
}

func (s *Scheduler) publishNotificationMessage(ctx context.Context) {
	err := s.producer.OpenChannel()
	if err != nil {
		s.log.Error("can't open channel", s.log.String("msg", err.Error()))
		return
	}
	defer func() {
		err := s.producer.CloseChannel()
		if err != nil {
			s.log.Error("can't close channel", s.log.String("msg", err.Error()))
		}
	}()

	from := time.Now()
	to := from.Add(s.interval)

	events, err := s.storage.EventListFilterByReminderIn(ctx, from.Unix(), to.Unix())
	if err != nil {
		s.log.Error("can't get events", s.log.String("msg", err.Error()))
		return
	}

	for _, e := range events {
		data, err := json.Marshal(NewMQEventNotification(e))
		if err != nil {
			s.log.Error("can't marshal event notification", s.log.String("msg", err.Error()))
			continue
		}
		err = s.producer.Publish(ctx, data)
		if err != nil {
			s.log.Error("can't publish event notification", s.log.String("msg", err.Error()))
		}
	}
}

func (s *Scheduler) clearEvents(ctx context.Context) {
	from := time.Now().AddDate(-1, 0, 0)
	to := from.Add(s.interval)

	events, err := s.storage.EventListFilterByStartDate(ctx, from.Unix(), to.Unix())
	if err != nil {
		s.log.Error("can't get events", s.log.String("msg", err.Error()))
		return
	}

	for _, e := range events {
		err := s.storage.RemoveEvent(ctx, e.ID)
		if err != nil {
			s.log.Error(
				"can't remove event",
				s.log.String("id", e.ID),
				s.log.String("msg", err.Error()),
			)
		}
	}
}

func startWorker(ctx context.Context, done chan struct{}, interval time.Duration, fn func()) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ctx.Done():
			close(done)
			return
		case <-ticker.C:
			fn()
		}
	}
}
