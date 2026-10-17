package analyzer

import (
	"slices"

	"github.com/Lazynx/spent/internal/statement"
)

const anomalyMedianFactor = 3

func detectAnomalies(txs []statement.Transaction, c Categorizer) []Anomaly {
	categories := make([]Category, len(txs))
	spendByCategory := make(map[Category][]int64)
	for i, t := range txs {
		if !t.IsExpense() {
			continue
		}
		categories[i] = c.Categorize(t)
		spendByCategory[categories[i]] = append(spendByCategory[categories[i]], t.Spent())
	}

	medians := make(map[Category]int64, len(spendByCategory))
	for cat, values := range spendByCategory {
		medians[cat] = median(values)
	}

	var out []Anomaly
	for i, t := range txs {
		if !t.IsExpense() {
			continue
		}
		cat := categories[i]
		if t.Spent() > anomalyMedianFactor*medians[cat] {
			out = append(out, Anomaly{Transaction: t, Category: cat, Median: medians[cat]})
		}
	}
	slices.SortStableFunc(out, func(a, b Anomaly) int {
		return a.Transaction.Date.Compare(b.Transaction.Date)
	})
	return out
}

func median(values []int64) int64 {
	if len(values) == 0 {
		return 0
	}
	sorted := slices.Clone(values)
	slices.Sort(sorted)
	mid := len(sorted) / 2
	if len(sorted)%2 == 1 {
		return sorted[mid]
	}
	return (sorted[mid-1] + sorted[mid]) / 2
}
