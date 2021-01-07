package memorystorage

import (
	"context"
	"sync"

	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/storage"
)

type EventDataStore struct {
	mu     sync.RWMutex
	events map[string]*app.Event
}

func New() *EventDataStore {
	return &EventDataStore{
		events: make(map[string]*app.Event),
	}
}

func (s *EventDataStore) NewEvent(ctx context.Context, e app.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.events[e.ID] != nil {
		return storage.ErrEventAlreadyExist
	}

	s.events[e.ID] = &e
	return nil
}

func (s *EventDataStore) UpdateEvent(ctx context.Context, e app.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.events[e.ID] == nil {
		return storage.ErrEventDoesNotExist
	}

	s.events[e.ID] = &e
	return nil
}

func (s *EventDataStore) RemoveEvent(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.events[id] == nil {
		return storage.ErrEventDoesNotExist
	}

	delete(s.events, id)
	return nil
}

func (s *EventDataStore) EventListFilterByStartDate(ctx context.Context, from int64, to int64) ([]app.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var events []app.Event

	for _, e := range s.events {
		if e.StartDate >= from && e.StartDate <= to {
			events = append(events, *e)
		}
	}

	if len(events) == 0 {
		return nil, storage.ErrNoEvents
	}
	return events, nil
}

func (s *EventDataStore) EventListFilterByReminderIn(ctx context.Context, from int64, to int64) ([]app.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var events []app.Event

	for _, e := range s.events {
		if e.RemindIn >= from && e.RemindIn <= to {
			events = append(events, *e)
		}
	}

	if len(events) == 0 {
		return nil, storage.ErrNoEvents
	}
	return events, nil
}
