package analyzer

import (
	"cmp"
	"slices"
	"strings"

	"github.com/Lazynx/spent/internal/statement"
)

const monthLayout = "2006-01"

func monthlyTotals(txs []statement.Transaction) []MonthTotal {
	byMonth := make(map[string]*MonthTotal)
	for _, t := range txs {
		key := t.Date.Format(monthLayout)
		m, ok := byMonth[key]
		if !ok {
			m = &MonthTotal{Month: key}
			byMonth[key] = m
		}
		if t.IsExpense() {
			m.Expenses += t.Spent()
		} else {
			m.Income += t.Amount
		}
	}

	out := make([]MonthTotal, 0, len(byMonth))
	for _, m := range byMonth {
		m.Net = m.Income - m.Expenses
		out = append(out, *m)
	}
	slices.SortFunc(out, func(a, b MonthTotal) int { return strings.Compare(a.Month, b.Month) })
	return out
}

func categoryTotals(txs []statement.Transaction, c Categorizer) []CategoryTotal {
	byCategory := make(map[Category]int64)
	for _, t := range txs {
		if t.IsExpense() {
			byCategory[c.Categorize(t)] += t.Spent()
		}
	}

	out := make([]CategoryTotal, 0, len(byCategory))
	for cat, total := range byCategory {
		out = append(out, CategoryTotal{Category: cat, Total: total})
	}
	slices.SortFunc(out, func(a, b CategoryTotal) int {
		return cmp.Or(cmp.Compare(b.Total, a.Total), cmp.Compare(a.Category, b.Category))
	})
	return out
}
