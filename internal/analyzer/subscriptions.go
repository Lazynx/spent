package analyzer

import "github.com/Lazynx/spent/internal/statement"

const (
	minSubscriptionOccurrences  = 3
	subscriptionAmountTolerance = 5
	minSubscriptionIntervalDays = 25
	maxSubscriptionIntervalDays = 35
)

func detectSubscriptions(txs []statement.Transaction) []Subscription {
	return nil
}
