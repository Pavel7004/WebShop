package domain

var (
	ErrItemNotFound          = NewError(404, "item_not_found", "Item not found")
	ErrInvalidId             = NewError(404, "invalid_id", "Can't parse id string")
	ErrUserNotFound          = NewError(404, "user_not_found", "User not found")
	ErrNoUpdate              = NewError(404, "update_not_specified", "There are no updates")
	ErrOrderNotProcessed     = NewError(404, "order_not_processed", "Order not processed")
	ErrOrderNotFound         = NewError(404, "order_not_found", "Order not found")
	ErrNoOrder               = NewError(404, "order_not_provided", "Order is nil")
	ErrOrderNotPaid          = NewError(404, "order_not_paid", "Order isn't paid")
	ErrOrderAlreadyDelivered = NewError(404, "order_delivered", "Order already delivered")
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
