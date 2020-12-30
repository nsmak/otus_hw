package storage

import "github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/app"

var (
	ErrEventAlreadyExist = NewError("event with this id already exist", nil)
	ErrEventDoesNotExist = NewError("event does not exist", nil)
	ErrNoEvents          = NewError("no one event", nil)
)

type Error struct {
	app.BaseError
}

func NewError(msg string, err error) *Error {
	return &Error{BaseError: app.BaseError{Message: msg, Err: err}}
}
