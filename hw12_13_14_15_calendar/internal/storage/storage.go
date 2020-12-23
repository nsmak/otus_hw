package storage

var (
	ErrEventAlreadyExist = &Error{Message: "event with this id already exist", Err: nil}
	ErrEventDoesNotExist = &Error{Message: "event does not exist", Err: nil}
	ErrNoEvents          = &Error{Message: "no one event", Err: nil}
)

type Error struct {
	Message string `json:"message"`
	Err     error  `json:"err,omitempty"`
}

func (e *Error) Error() string {
	if e.Err != nil {
		e.Message = e.Message + " --> " + e.Err.Error()
	}
	return e.Message
}
func (e *Error) Unwrap() error {
	return e.Err
}
