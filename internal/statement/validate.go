package statement

import (
	"fmt"
	"strings"
	"time"
)

var supportedCurrencies = map[string]struct{}{
	"KZT": {},
	"USD": {},
	"EUR": {},
}

func (s Statement) Validate(now time.Time) error {
	if len(s.Transactions) == 0 {
		return ErrEmptyStatement
	}
	for i, t := range s.Transactions {
		if err := t.validate(now); err != nil {
			return fmt.Errorf("transaction %d: %w", i, err)
		}
	}
	return nil
}

func (t Transaction) validate(now time.Time) error {
	if t.Date.IsZero() {
		return ErrZeroDate
	}
	if t.Date.After(now) {
		return ErrFutureDate
	}
	if t.Amount == 0 {
		return ErrZeroAmount
	}
	if _, ok := supportedCurrencies[t.Currency]; !ok {
		return ErrInvalidCurrency
	}
	if strings.TrimSpace(t.Description) == "" {
		return ErrEmptyDescription
	}
	return nil
}
