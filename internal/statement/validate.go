package statement

import "time"

var supportedCurrencies = map[string]struct{}{
	"KZT": {},
	"USD": {},
	"EUR": {},
}

func (s Statement) Validate(now time.Time) error {
	return errNotImplemented
}
