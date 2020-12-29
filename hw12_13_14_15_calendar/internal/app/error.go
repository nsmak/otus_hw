package app

type BaseError struct {
	Message string `json:"message"`
	Err     error  `json:"err,omitempty"`
}

func (e *BaseError) Error() string {
	if e.Err != nil {
		e.Message = e.Message + " --> " + e.Err.Error()
	}
	return e.Message
}
func (e *BaseError) Unwrap() error {
	return e.Err
}
