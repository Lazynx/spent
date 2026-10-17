package analyzer

import (
	"errors"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/Lazynx/spent/internal/statement"
)

func tx(date string, amount int64, desc, mcc string) statement.Transaction {
	d, err := time.Parse("2006-01-02", date)
	if err != nil {
		panic(err)
	}
	return statement.Transaction{Date: d, Amount: amount, Currency: "KZT", Description: desc, MCC: mcc}
}

func loadSample(t *testing.T) statement.Statement {
	t.Helper()
	f, err := os.Open("../../testdata/sample.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	st, err := statement.CSVParser{}.Parse(f)
	if err != nil {
		t.Fatal(err)
	}
	return st
}

func TestAnalyzeSample(t *testing.T) {
	rep, err := Analyze(loadSample(t), NewRuleCategorizer())
	if err != nil {
		t.Fatalf("Analyze() error = %v", err)
	}

	wantMonthly := []MonthTotal{
		{"2026-01", 50000000, 2550000, 47450000},
		{"2026-02", 50000000, 2950000, 47050000},
		{"2026-03", 50000000, 10850000, 39150000},
	}
	if !reflect.DeepEqual(rep.Monthly, wantMonthly) {
		t.Errorf("Monthly = %+v, want %+v", rep.Monthly, wantMonthly)
	}

	wantCategories := []CategoryTotal{
		{Food, 15600000},
		{Entertainment, 750000},
	}
	if !reflect.DeepEqual(rep.ByCategory, wantCategories) {
		t.Errorf("ByCategory = %+v, want %+v", rep.ByCategory, wantCategories)
	}

	wantMerchants := []MerchantTotal{
		{"magnum", 13900000},
		{"glovo", 1700000},
		{"netflix.com", 750000},
	}
	if !reflect.DeepEqual(rep.TopMerchants, wantMerchants) {
		t.Errorf("TopMerchants = %+v, want %+v", rep.TopMerchants, wantMerchants)
	}

	wantSubs := []Subscription{
		{"netflix.com", 250000, 3, time.Date(2026, 3, 13, 0, 0, 0, 0, time.UTC)},
	}
	if !reflect.DeepEqual(rep.Subscriptions, wantSubs) {
		t.Errorf("Subscriptions = %+v, want %+v", rep.Subscriptions, wantSubs)
	}

	if len(rep.Anomalies) != 1 {
		t.Fatalf("Anomalies = %+v, want 1", rep.Anomalies)
	}
	a := rep.Anomalies[0]
	if a.Transaction.Amount != -9000000 || a.Category != Food || a.Median != 1550000 {
		t.Errorf("anomaly = %+v", a)
	}
}

func TestCheckBalance(t *testing.T) {
	tests := []struct {
		name    string
		opening int64
		closing int64
		wantErr bool
	}{
		{"matches", 1000, 600, false},
		{"mismatch", 1000, 601, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := statement.Statement{
				OpeningBalance: tt.opening,
				ClosingBalance: tt.closing,
				Transactions: []statement.Transaction{
					tx("2026-01-01", -500, "Shop", ""),
					tx("2026-01-02", 100, "Refund", ""),
				},
			}
			err := checkBalance(st)
			if tt.wantErr != errors.Is(err, ErrBalanceMismatch) {
				t.Fatalf("checkBalance() = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMedian(t *testing.T) {
	tests := []struct {
		name string
		in   []int64
		want int64
	}{
		{"empty", nil, 0},
		{"single", []int64{7}, 7},
		{"odd", []int64{9, 1, 5}, 5},
		{"even", []int64{4, 1, 3, 2}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := median(tt.in); got != tt.want {
				t.Errorf("median(%v) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

func TestMedianDoesNotMutateInput(t *testing.T) {
	in := []int64{3, 1, 2}
	median(in)
	if !reflect.DeepEqual(in, []int64{3, 1, 2}) {
		t.Errorf("input mutated: %v", in)
	}
}

func TestNormalize(t *testing.T) {
	if got := normalize("  NetFlix   .COM "); got != "netflix .com" {
		t.Errorf("normalize = %q", got)
	}
}

func TestDetectSubscriptions(t *testing.T) {
	tests := []struct {
		name string
		txs  []statement.Transaction
		want int
	}{
		{
			"monthly with small price drift",
			[]statement.Transaction{
				tx("2026-01-10", -100000, "Spotify", ""),
				tx("2026-02-09", -104000, "SPOTIFY", ""),
				tx("2026-03-11", -100000, "spotify", ""),
			},
			1,
		},
		{
			"only two occurrences",
			[]statement.Transaction{
				tx("2026-01-10", -100000, "Spotify", ""),
				tx("2026-02-09", -100000, "Spotify", ""),
			},
			0,
		},
		{
			"amount drift above 5 percent",
			[]statement.Transaction{
				tx("2026-01-10", -100000, "Spotify", ""),
				tx("2026-02-09", -106000, "Spotify", ""),
				tx("2026-03-11", -100000, "Spotify", ""),
			},
			0,
		},
		{
			"interval too short",
			[]statement.Transaction{
				tx("2026-01-10", -100000, "Gym", ""),
				tx("2026-01-20", -100000, "Gym", ""),
				tx("2026-01-30", -100000, "Gym", ""),
			},
			0,
		},
		{
			"interval too long",
			[]statement.Transaction{
				tx("2026-01-10", -100000, "Gym", ""),
				tx("2026-03-10", -100000, "Gym", ""),
				tx("2026-05-10", -100000, "Gym", ""),
			},
			0,
		},
		{
			"income is ignored",
			[]statement.Transaction{
				tx("2026-01-01", 100000, "Salary", ""),
				tx("2026-02-01", 100000, "Salary", ""),
				tx("2026-03-01", 100000, "Salary", ""),
			},
			0,
		},
		{
			"unsorted input",
			[]statement.Transaction{
				tx("2026-03-11", -100000, "Spotify", ""),
				tx("2026-01-10", -100000, "Spotify", ""),
				tx("2026-02-09", -100000, "Spotify", ""),
			},
			1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectSubscriptions(tt.txs)
			if len(got) != tt.want {
				t.Fatalf("detectSubscriptions() = %+v, want %d", got, tt.want)
			}
		})
	}
}

func TestDetectAnomaliesThresholdIsStrict(t *testing.T) {
	c := NewRuleCategorizer()
	txs := []statement.Transaction{
		tx("2026-01-01", -100, "Magnum", ""),
		tx("2026-01-02", -100, "Magnum", ""),
		tx("2026-01-03", -300, "Magnum", ""),
	}
	if got := detectAnomalies(txs, c); len(got) != 0 {
		t.Errorf("exactly 3x median flagged: %+v", got)
	}

	txs[2].Amount = -301
	if got := detectAnomalies(txs, c); len(got) != 1 {
		t.Errorf("above 3x median not flagged: %+v", got)
	}
}

func TestTopMerchantsLimit(t *testing.T) {
	txs := []statement.Transaction{
		tx("2026-01-01", -300, "A", ""),
		tx("2026-01-02", -200, "B", ""),
		tx("2026-01-03", -100, "C", ""),
	}
	got := topMerchants(txs, 2)
	want := []MerchantTotal{{"a", 300}, {"b", 200}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("topMerchants() = %+v, want %+v", got, want)
	}
}

func TestRuleCategorizer(t *testing.T) {
	c := NewRuleCategorizer()
	tests := []struct {
		name string
		tx   statement.Transaction
		want Category
	}{
		{"mcc", tx("2026-01-01", -1, "Unknown shop", "5411"), Food},
		{"mcc wins over keyword", tx("2026-01-01", -1, "Netflix", "5411"), Food},
		{"keyword", tx("2026-01-01", -1, "NETFLIX.COM 123", ""), Entertainment},
		{"unknown mcc falls to keyword", tx("2026-01-01", -1, "Uber trip", "0000"), Transport},
		{"fallback", tx("2026-01-01", -1, "Random", ""), Other},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.Categorize(tt.tx); got != tt.want {
				t.Errorf("Categorize() = %s, want %s", got, tt.want)
			}
		})
	}
}
