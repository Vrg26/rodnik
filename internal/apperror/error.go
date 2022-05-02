package apperror

const (
	NoType = ErrorType(iota)
	NotFound
	BadRequest
	Conflict
	Internal
	Authorization
	PaymentRequired
)

type ErrorType uint

type AppError struct {
	Type    ErrorType `json:"-"`
	Message string    `json:"message,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func (t ErrorType) New(message string) *AppError {
	return &AppError{
		Type:    t,
		Message: message,
	}
}
