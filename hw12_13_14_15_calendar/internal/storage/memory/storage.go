package memorystorage

import (
	"context"
	"sync"

	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/storage"
)

var (
	ErrEventAlreadyExist = &MemDBError{Message: "event with this id already exist", Err: nil}
	ErrEventDoesNotExist = &MemDBError{Message: "event does not exist", Err: nil}
	ErrNoEvents          = &MemDBError{Message: "no one event", Err: nil}
)

type MemDBError struct {
	Message string `json:"message"`
	Err     error  `json:"err,omitempty"`
}

func (e *MemDBError) Error() string {
	if e.Err != nil {
		e.Message = e.Message + " --> " + e.Err.Error()
	}
	return e.Message
}
func (e *MemDBError) Unwrap() error {
	return e.Err
}

type Storage struct {
	mu     sync.RWMutex
	events map[string]*storage.Event
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) NewEvent(ctx context.Context, e storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.events[e.ID] != nil {
		return ErrEventAlreadyExist
	}

	s.events[e.ID] = &e
	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, e storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.events[e.ID] == nil {
		return ErrEventDoesNotExist
	}

	s.events[e.ID] = &e
	return nil
}

func (s *Storage) RemoveEvent(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.events[id] == nil {
		return ErrEventDoesNotExist
	}

	delete(s.events, id)
	return nil
}

func (s *Storage) EventList(ctx context.Context, from int64, to int64) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var events []storage.Event

	for _, e := range s.events {
		if e.StartDate >= from && e.StartDate <= to {
			events = append(events, *e)
		}
	}

	if len(events) == 0 {
		return nil, ErrNoEvents
	}
	return events, nil
}
