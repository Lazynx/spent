package statement

import (
	"errors"
	"os"
	"strings"
	"testing"
	"time"
)

var now = time.Date(2026, 10, 17, 12, 0, 0, 0, time.UTC)

func validTx() Transaction {
	return Transaction{
		Date:        time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC),
		Amount:      -150000,
		Currency:    "KZT",
		Description: "Magnum",
		MCC:         "5411",
	}
}

func TestValidate(t *testing.T) {
	with := func(mod func(*Transaction)) Statement {
		tx := validTx()
		mod(&tx)
		return Statement{Transactions: []Transaction{tx}}
	}

	tests := []struct {
		name string
		st   Statement
		want error
	}{
		{"valid", Statement{Transactions: []Transaction{validTx()}}, nil},
		{"no transactions", Statement{}, ErrEmptyStatement},
		{"zero date", with(func(tx *Transaction) { tx.Date = time.Time{} }), ErrZeroDate},
		{"future date", with(func(tx *Transaction) { tx.Date = now.Add(time.Hour) }), ErrFutureDate},
		{"zero amount", with(func(tx *Transaction) { tx.Amount = 0 }), ErrZeroAmount},
		{"empty currency", with(func(tx *Transaction) { tx.Currency = "" }), ErrInvalidCurrency},
		{"unsupported currency", with(func(tx *Transaction) { tx.Currency = "XYZ" }), ErrInvalidCurrency},
		{"empty description", with(func(tx *Transaction) { tx.Description = "" }), ErrEmptyDescription},
		{"blank description", with(func(tx *Transaction) { tx.Description = "   " }), ErrEmptyDescription},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.st.Validate(now)
			if !errors.Is(err, tt.want) {
				t.Fatalf("Validate() = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestValidateReportsTransactionIndex(t *testing.T) {
	bad := validTx()
	bad.Amount = 0
	st := Statement{Transactions: []Transaction{validTx(), bad}}

	err := st.Validate(now)
	if err == nil || !strings.Contains(err.Error(), "transaction 1") {
		t.Fatalf("Validate() = %v, want error mentioning transaction 1", err)
	}
}

func TestCSVParserSample(t *testing.T) {
	f, err := os.Open("../../testdata/sample.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	st, err := CSVParser{}.Parse(f)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if st.OpeningBalance != 100000000 || st.ClosingBalance != 233650000 {
		t.Errorf("balances = %d / %d", st.OpeningBalance, st.ClosingBalance)
	}
	if len(st.Transactions) != 12 {
		t.Fatalf("got %d transactions, want 12", len(st.Transactions))
	}

	first := st.Transactions[0]
	want := Transaction{
		Date:        time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		Amount:      50000000,
		Currency:    "KZT",
		Description: "Salary ACME",
	}
	if first != want {
		t.Errorf("first = %+v, want %+v", first, want)
	}
	if st.Transactions[1].MCC != "5411" {
		t.Errorf("MCC = %q, want 5411", st.Transactions[1].MCC)
	}
}

func TestCSVParserErrors(t *testing.T) {
	const balances = "opening_balance,0\nclosing_balance,0\n"
	tests := []struct {
		name string
		in   string
	}{
		{"missing balances", "date,amount,currency,description,mcc\n2026-01-01,-1,KZT,Shop,\n"},
		{"missing closing", "opening_balance,0\n2026-01-01,-1,KZT,Shop,\n"},
		{"bad balance", "opening_balance,abc\nclosing_balance,0\n"},
		{"bad date", balances + "01/01/2026,-1,KZT,Shop,\n"},
		{"bad amount", balances + "2026-01-01,1.5,KZT,Shop,\n"},
		{"wrong field count", balances + "2026-01-01,-1,KZT\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CSVParser{}.Parse(strings.NewReader(tt.in))
			if !errors.Is(err, ErrInvalidFormat) {
				t.Fatalf("Parse() = %v, want ErrInvalidFormat", err)
			}
		})
	}
}

func TestCSVParserNormalizesCurrency(t *testing.T) {
	in := "opening_balance,0\nclosing_balance,-1\n2026-01-01,-1,kzt, Shop ,\n"
	st, err := CSVParser{}.Parse(strings.NewReader(in))
	if err != nil {
		t.Fatal(err)
	}
	got := st.Transactions[0]
	if got.Currency != "KZT" || got.Description != "Shop" {
		t.Errorf("got currency %q description %q", got.Currency, got.Description)
	}
}
