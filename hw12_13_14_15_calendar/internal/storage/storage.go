package storage

type EventDataStore interface {
	NewEvent(e Event) error
	UpdateEvent(e Event) error
	RemoveEvent(id string) error
	EventList(from int64, to int64) ([]Event, error)
}
