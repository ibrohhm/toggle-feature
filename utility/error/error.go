package error

type Error struct {
	Message    string `json:"message"`
	HttpStatus int    `json:"http_status"`
}

func (e Error) Error() string {
	return e.Message
}

func New(message string, httpStatus int) Error {
	return Error{
		Message:    message,
		HttpStatus: httpStatus,
	}
}
