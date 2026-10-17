package analyzer

import (
	"fmt"

	"github.com/Lazynx/spent/internal/statement"
)

const topMerchantsLimit = 5

func Analyze(s statement.Statement, c Categorizer) (Report, error) {
	if err := checkBalance(s); err != nil {
		return Report{}, err
	}
	return Report{
		Monthly:       monthlyTotals(s.Transactions),
		ByCategory:    categoryTotals(s.Transactions, c),
		TopMerchants:  topMerchants(s.Transactions, topMerchantsLimit),
		Subscriptions: detectSubscriptions(s.Transactions),
		Anomalies:     detectAnomalies(s.Transactions, c),
	}, nil
}

func checkBalance(s statement.Statement) error {
	var sum int64
	for _, t := range s.Transactions {
		sum += t.Amount
	}
	if got := s.OpeningBalance + sum; got != s.ClosingBalance {
		return fmt.Errorf("opening %d + transactions %d = %d, but closing is %d: %w",
			s.OpeningBalance, sum, got, s.ClosingBalance, ErrBalanceMismatch)
	}
	return nil
}
