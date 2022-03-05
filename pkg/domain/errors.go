package domain

var (
	ErrItemNotFound = NewError(404, "item_not_found", "Item not found")
	ErrInvalidId    = NewError(404, "invalid_id", "Can't parse id string")
)

type Error struct {
	CodeHTTP int    `json:"-"`
	Code     string `json:"code,omitempty"`
	Message  string `json:"message"`
}

func (err *Error) Error() string {
	return err.Message
}

func NewError(httpCode int, code, message string) error {
	return &Error{
		CodeHTTP: httpCode,
		Code:     code,
		Message:  message,
	}
}
