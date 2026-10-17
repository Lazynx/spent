package main

import (
	"fmt"
	"io"

	"github.com/Lazynx/spent/internal/analyzer"
)

func printReport(w io.Writer, r analyzer.Report) error {
	var err error
	p := func(format string, args ...any) {
		if err == nil {
			_, err = fmt.Fprintf(w, format, args...)
		}
	}

	p("Monthly totals\n")
	for _, m := range r.Monthly {
		p("  %s  income %12s  expenses %12s  net %12s\n",
			m.Month, formatMoney(m.Income), formatMoney(m.Expenses), formatMoney(m.Net))
	}

	p("\nSpending by category\n")
	for _, c := range r.ByCategory {
		p("  %-14s %12s\n", c.Category, formatMoney(c.Total))
	}

	p("\nTop merchants\n")
	for _, m := range r.TopMerchants {
		p("  %-24s %12s\n", m.Merchant, formatMoney(m.Total))
	}

	p("\nSubscriptions\n")
	if len(r.Subscriptions) == 0 {
		p("  none\n")
	}
	for _, s := range r.Subscriptions {
		p("  %-24s %12s  x%d, last %s\n",
			s.Merchant, formatMoney(s.Amount), s.Occurrences, s.LastDate.Format("2006-01-02"))
	}

	p("\nAnomalies\n")
	if len(r.Anomalies) == 0 {
		p("  none\n")
	}
	for _, a := range r.Anomalies {
		p("  %s  %-24s %12s  (%s median %s)\n",
			a.Transaction.Date.Format("2006-01-02"), a.Transaction.Description,
			formatMoney(a.Transaction.Spent()), a.Category, formatMoney(a.Median))
	}

	return err
}

func formatMoney(tiyn int64) string {
	sign := ""
	if tiyn < 0 {
		sign = "-"
		tiyn = -tiyn
	}
	return fmt.Sprintf("%s%d.%02d", sign, tiyn/100, tiyn%100)
}
