package statement

import "time"

type Transaction struct {
	Date        time.Time
	Amount      int64
	Currency    string
	Description string
	MCC         string
}

type Statement struct {
	OpeningBalance int64
	ClosingBalance int64
	Transactions   []Transaction
}
