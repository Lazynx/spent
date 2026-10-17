package analyzer

import (
	"cmp"
	"slices"
	"strings"

	"github.com/Lazynx/spent/internal/statement"
)

func normalize(description string) string {
	return strings.Join(strings.Fields(strings.ToLower(description)), " ")
}

func topMerchants(txs []statement.Transaction, limit int) []MerchantTotal {
	byMerchant := make(map[string]int64)
	for _, t := range txs {
		if t.IsExpense() {
			byMerchant[normalize(t.Description)] += t.Spent()
		}
	}

	out := make([]MerchantTotal, 0, len(byMerchant))
	for m, total := range byMerchant {
		out = append(out, MerchantTotal{Merchant: m, Total: total})
	}
	slices.SortFunc(out, func(a, b MerchantTotal) int {
		return cmp.Or(cmp.Compare(b.Total, a.Total), cmp.Compare(a.Merchant, b.Merchant))
	})
	if len(out) > limit {
		out = out[:limit]
	}
	return out
}
