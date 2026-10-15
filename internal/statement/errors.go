package statement

import "errors"

var (
	ErrEmptyStatement   = errors.New("statement has no transactions")
	ErrZeroDate         = errors.New("transaction date is zero")
	ErrFutureDate       = errors.New("transaction date is in the future")
	ErrZeroAmount       = errors.New("transaction amount is zero")
	ErrInvalidCurrency  = errors.New("currency is empty or unsupported")
	ErrEmptyDescription = errors.New("transaction description is empty")
	ErrInvalidFormat    = errors.New("invalid statement format")

	errNotImplemented = errors.New("not implemented")
)
