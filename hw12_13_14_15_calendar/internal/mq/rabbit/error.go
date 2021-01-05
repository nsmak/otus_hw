package rabbit

import "github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/app"

type Error struct {
	app.BaseError
}

func (e *Error) Error() string {
	if e.Err != nil {
		e.Message = "[rmq] " + e.Message + " --> " + e.Err.Error()
	}
	return e.Message
}

func NewError(msg string, err error) *Error {
	return &Error{BaseError: app.BaseError{Message: msg, Err: err}}
}

var ErrChannelIsNil = NewError("channel is nil", nil)
