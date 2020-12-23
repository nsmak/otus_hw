package app

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type ProcessingError struct {
	Message string `json:"message"`
	Err     error  `json:"err,omitempty"`
}

func (e *ProcessingError) Error() string {
	if e.Err != nil {
		e.Message = e.Message + " --> " + e.Err.Error()
	}
	return e.Message
}
func (e *ProcessingError) Unwrap() error {
	return e.Err
}

type Logger interface {
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	String(key string, val string) zap.Field
	Int64(key string, val int64) zap.Field
	Duration(key string, val time.Duration) zap.Field
}

//go:generate mockgen -destination=./mock_storage_test.go -package=app_test . Storage
type Storage interface {
	NewEvent(ctx context.Context, e Event) error
	UpdateEvent(ctx context.Context, e Event) error
	RemoveEvent(ctx context.Context, id string) error
	EventList(ctx context.Context, from int64, to int64) ([]Event, error)
}

type App struct {
	log     Logger
	storage Storage
}

func New(logger Logger, storage Storage) *App {
	return &App{log: logger, storage: storage}
}

func (a *App) CreateEvent(ctx context.Context, e Event) error {
	a.log.Info("create event")
	err := a.storage.NewEvent(ctx, e)
	if err != nil {
		return &ProcessingError{
			Message: "can't create event",
			Err:     err,
		}
	}
	return nil
}

func (a *App) UpdateEvent(ctx context.Context, e Event) error {
	err := a.storage.UpdateEvent(ctx, e)
	if err != nil {
		return &ProcessingError{
			Message: "can't update event",
			Err:     err,
		}
	}
	return nil
}

func (a *App) RemoveEvent(ctx context.Context, id string) error {
	err := a.storage.RemoveEvent(ctx, id)
	if err != nil {
		return &ProcessingError{
			Message: "can't remove event",
			Err:     err,
		}
	}
	return nil
}

func (a *App) Events(ctx context.Context, from int64, to int64) ([]Event, error) {
	events, err := a.storage.EventList(ctx, from, to)
	if err != nil {
		return nil, &ProcessingError{
			Message: "can't get events",
			Err:     err,
		}
	}
	return events, nil
}
