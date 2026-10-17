package statement

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

const (
	dateLayout        = "2006-01-02"
	transactionFields = 5
)

type CSVParser struct{}

func (CSVParser) Parse(r io.Reader) (Statement, error) {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	cr.TrimLeadingSpace = true

	var (
		st         Statement
		hasOpening bool
		hasClosing bool
	)

	for {
		rec, err := cr.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return Statement{}, fmt.Errorf("read csv: %w: %w", ErrInvalidFormat, err)
		}
		line, _ := cr.FieldPos(0)

		switch strings.ToLower(strings.TrimSpace(rec[0])) {
		case "opening_balance":
			v, err := parseBalance(rec)
			if err != nil {
				return Statement{}, fmt.Errorf("line %d: %w", line, err)
			}
			st.OpeningBalance, hasOpening = v, true
		case "closing_balance":
			v, err := parseBalance(rec)
			if err != nil {
				return Statement{}, fmt.Errorf("line %d: %w", line, err)
			}
			st.ClosingBalance, hasClosing = v, true
		case "date":
		default:
			t, err := parseTransaction(rec)
			if err != nil {
				return Statement{}, fmt.Errorf("line %d: %w", line, err)
			}
			st.Transactions = append(st.Transactions, t)
		}
	}

	if !hasOpening || !hasClosing {
		return Statement{}, fmt.Errorf("opening_balance and closing_balance are required: %w", ErrInvalidFormat)
	}
	return st, nil
}

func parseBalance(rec []string) (int64, error) {
	if len(rec) != 2 {
		return 0, fmt.Errorf("balance row needs 2 fields, got %d: %w", len(rec), ErrInvalidFormat)
	}
	v, err := strconv.ParseInt(strings.TrimSpace(rec[1]), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("balance %q: %w", rec[1], ErrInvalidFormat)
	}
	return v, nil
}

func parseTransaction(rec []string) (Transaction, error) {
	if len(rec) != transactionFields {
		return Transaction{}, fmt.Errorf("transaction row needs %d fields, got %d: %w", transactionFields, len(rec), ErrInvalidFormat)
	}

	date, err := time.Parse(dateLayout, strings.TrimSpace(rec[0]))
	if err != nil {
		return Transaction{}, fmt.Errorf("date %q: %w", rec[0], ErrInvalidFormat)
	}
	amount, err := strconv.ParseInt(strings.TrimSpace(rec[1]), 10, 64)
	if err != nil {
		return Transaction{}, fmt.Errorf("amount %q: %w", rec[1], ErrInvalidFormat)
	}

	return Transaction{
		Date:        date,
		Amount:      amount,
		Currency:    strings.ToUpper(strings.TrimSpace(rec[2])),
		Description: strings.TrimSpace(rec[3]),
		MCC:         strings.TrimSpace(rec[4]),
	}, nil
}
