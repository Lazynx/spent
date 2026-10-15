package analyzer

import (
	"time"

	"github.com/Lazynx/spent/internal/statement"
)

type Report struct {
	Monthly       []MonthTotal
	ByCategory    []CategoryTotal
	TopMerchants  []MerchantTotal
	Subscriptions []Subscription
	Anomalies     []Anomaly
}

type MonthTotal struct {
	Month    string
	Income   int64
	Expenses int64
	Net      int64
}

type CategoryTotal struct {
	Category Category
	Total    int64
}

type MerchantTotal struct {
	Merchant string
	Total    int64
}

type Subscription struct {
	Merchant    string
	Amount      int64
	Occurrences int
	LastDate    time.Time
}

type Anomaly struct {
	Transaction statement.Transaction
	Category    Category
	Median      int64
}
