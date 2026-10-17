package analyzer

import (
	"cmp"
	"slices"
	"time"

	"github.com/Lazynx/spent/internal/statement"
)

const (
	minSubscriptionOccurrences  = 3
	subscriptionAmountTolerance = 5
	minSubscriptionIntervalDays = 25
	maxSubscriptionIntervalDays = 35
)

func detectSubscriptions(txs []statement.Transaction) []Subscription {
	groups := make(map[string][]statement.Transaction)
	for _, t := range txs {
		if t.IsExpense() {
			key := normalize(t.Description)
			groups[key] = append(groups[key], t)
		}
	}

	var subs []Subscription
	for merchant, g := range groups {
		if len(g) < minSubscriptionOccurrences {
			continue
		}
		slices.SortStableFunc(g, func(a, b statement.Transaction) int { return a.Date.Compare(b.Date) })
		if !isRecurring(g) {
			continue
		}
		last := g[len(g)-1]
		subs = append(subs, Subscription{
			Merchant:    merchant,
			Amount:      last.Spent(),
			Occurrences: len(g),
			LastDate:    last.Date,
		})
	}
	slices.SortFunc(subs, func(a, b Subscription) int { return cmp.Compare(a.Merchant, b.Merchant) })
	return subs
}

func isRecurring(g []statement.Transaction) bool {
	ref := g[0].Spent()
	for i, t := range g {
		if !withinTolerance(t.Spent(), ref) {
			return false
		}
		if i == 0 {
			continue
		}
		days := daysBetween(g[i-1].Date, t.Date)
		if days < minSubscriptionIntervalDays || days > maxSubscriptionIntervalDays {
			return false
		}
	}
	return true
}

func withinTolerance(amount, ref int64) bool {
	diff := amount - ref
	if diff < 0 {
		diff = -diff
	}
	return diff*100 <= ref*subscriptionAmountTolerance
}

func daysBetween(a, b time.Time) int {
	return int(dayStart(b).Sub(dayStart(a)).Hours() / 24)
}

func dayStart(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}
