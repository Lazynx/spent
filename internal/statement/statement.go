package statement

import "time"

type Transaction struct {
	Date        time.Time
	Amount      int64
	Currency    string
	Description string
	MCC         string
}

func (t Transaction) IsExpense() bool { return t.Amount < 0 }

func (t Transaction) Spent() int64 { return -t.Amount }

type Statement struct {
	OpeningBalance int64
	ClosingBalance int64
	Transactions   []Transaction
}
