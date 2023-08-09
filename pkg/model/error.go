package model

type Error struct {
	Status  int
	Code    string
	Message string
}

func (e Error) Error() string {
	return e.Message
}

// NewError new a error with message
func NewError(status int, code, msg string) error {
	return Error{
		Status:  status,
		Code:    code,
		Message: msg,
	}
}
