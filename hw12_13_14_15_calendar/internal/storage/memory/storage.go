package memorystorage

import (
	"sync"

	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/storage"
)

var (
	ErrEventAlreadyExist = &memDBError{Message: "event with this id already exist", Err: nil}
	ErrEventDoesNotExist = &memDBError{Message: "event does not exist", Err: nil}
	ErrNoEvents          = &memDBError{Message: "no one event", Err: nil}
)

type memDBError struct {
	Message string `json:"message"`
	Err     error  `json:"err,omitempty"`
}

func (e *memDBError) Error() string {
	if e.Err != nil {
		e.Message = e.Message + " --> " + e.Err.Error()
	}
	return e.Message
}
func (e *memDBError) Unwrap() error {
	return e.Err
}

type Storage struct {
	mu     sync.RWMutex
	events map[string]*storage.Event
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) NewEvent(e storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.events[e.ID] != nil {
		return ErrEventAlreadyExist
	}

	s.events[e.ID] = &e
	return nil
}

func (s *Storage) UpdateEvent(e storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.events[e.ID] == nil {
		return ErrEventDoesNotExist
	}

	s.events[e.ID] = &e
	return nil
}

func (s *Storage) RemoveEvent(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.events[id] == nil {
		return ErrEventDoesNotExist
	}

	delete(s.events, id)
	return nil
}

func (s *Storage) EventList(from int64, to int64) ([]storage.Event, error) {
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
