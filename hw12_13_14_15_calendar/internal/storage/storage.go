package storage

import "context"

type EventDataStore interface {
	NewEvent(ctx context.Context, e Event) error
	UpdateEvent(ctx context.Context, e Event) error
	RemoveEvent(ctx context.Context, id string) error
	EventList(ctx context.Context, from int64, to int64) ([]Event, error)
}
