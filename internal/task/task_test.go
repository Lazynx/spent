package task

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Lazynx/spent/internal/analyzer"
	"github.com/Lazynx/spent/internal/statement"
)

var fixedNow = time.Date(2026, 10, 17, 12, 0, 0, 0, time.UTC)

func newTestService() Service {
	n := 0
	newID := func() string {
		n++
		return fmt.Sprintf("task-%d", n)
	}
	analyze := func(s statement.Statement) (analyzer.Report, error) {
		return analyzer.Analyze(s, analyzer.NewRuleCategorizer())
	}
	return NewService(NewMemoryStore(), analyze, func() time.Time { return fixedNow }, newID)
}

func sampleStatement(t *testing.T) statement.Statement {
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

func TestStatusTransitions(t *testing.T) {
	all := []Status{Pending, Processing, Completed, Failed}
	allowed := map[[2]Status]bool{
		{Pending, Processing}:   true,
		{Processing, Completed}: true,
		{Processing, Failed}:    true,
	}
	for _, from := range all {
		for _, to := range all {
			want := allowed[[2]Status{from, to}]
			if got := from.CanTransitionTo(to); got != want {
				t.Errorf("%s -> %s = %v, want %v", from, to, got, want)
			}
		}
	}
}

func TestStatusString(t *testing.T) {
	tests := map[Status]string{
		Pending: "pending", Processing: "processing",
		Completed: "completed", Failed: "failed", Status(99): "unknown",
	}
	for s, want := range tests {
		if got := s.String(); got != want {
			t.Errorf("Status(%d).String() = %q, want %q", int(s), got, want)
		}
	}
}

func TestCreate(t *testing.T) {
	svc := newTestService()
	st := sampleStatement(t)

	tk, err := svc.Create(st)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if tk.ID != "task-1" || tk.Status != Pending || tk.Report != nil || tk.FailReason != "" {
		t.Errorf("unexpected task: %+v", tk)
	}
	if len(tk.Statement.Transactions) != 12 {
		t.Errorf("transactions = %d, want 12", len(tk.Statement.Transactions))
	}
}

func TestCreateInvalid(t *testing.T) {
	tests := []struct {
		name string
		mod  func(*statement.Statement)
		want error
	}{
		{"empty", func(s *statement.Statement) { s.Transactions = nil }, statement.ErrEmptyStatement},
		{"zero date", func(s *statement.Statement) { s.Transactions[0].Date = time.Time{} }, statement.ErrZeroDate},
		{"future date", func(s *statement.Statement) { s.Transactions[0].Date = fixedNow.Add(24 * time.Hour) }, statement.ErrFutureDate},
		{"zero amount", func(s *statement.Statement) { s.Transactions[0].Amount = 0 }, statement.ErrZeroAmount},
		{"bad currency", func(s *statement.Statement) { s.Transactions[0].Currency = "" }, statement.ErrInvalidCurrency},
		{"empty description", func(s *statement.Statement) { s.Transactions[0].Description = "" }, statement.ErrEmptyDescription},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService()
			st := sampleStatement(t)
			tt.mod(&st)

			tk, err := svc.Create(st)
			if !errors.Is(err, tt.want) {
				t.Fatalf("Create() error = %v, want %v", err, tt.want)
			}
			if tk != nil {
				t.Errorf("task created despite invalid input: %+v", tk)
			}
			if _, err := svc.Get("task-1"); !errors.Is(err, ErrTaskNotFound) {
				t.Errorf("invalid input was stored: %v", err)
			}
		})
	}
}

func TestCreateCopiesTransactions(t *testing.T) {
	svc := newTestService()
	st := sampleStatement(t)

	tk, err := svc.Create(st)
	if err != nil {
		t.Fatal(err)
	}
	st.Transactions[0].Amount = 1
	if tk.Statement.Transactions[0].Amount == 1 {
		t.Error("task shares transactions slice with caller")
	}
}

func TestGet(t *testing.T) {
	svc := newTestService()
	created, err := svc.Create(sampleStatement(t))
	if err != nil {
		t.Fatal(err)
	}

	got, err := svc.Get(created.ID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("Get() ID = %q, want %q", got.ID, created.ID)
	}

	if _, err := svc.Get("missing"); !errors.Is(err, ErrTaskNotFound) {
		t.Errorf("Get(missing) = %v, want ErrTaskNotFound", err)
	}
}

func TestProcessSuccess(t *testing.T) {
	svc := newTestService()
	created, err := svc.Create(sampleStatement(t))
	if err != nil {
		t.Fatal(err)
	}

	got, err := svc.Process(created.ID)
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	if got.Status != Completed {
		t.Fatalf("status = %s, want completed", got.Status)
	}
	if got.Report == nil {
		t.Fatal("report is nil")
	}
	if len(got.Report.Monthly) != 3 || len(got.Report.Subscriptions) != 1 || len(got.Report.Anomalies) != 1 {
		t.Errorf("unexpected report: %+v", got.Report)
	}
	if got.Report.Monthly[2].Expenses != 10850000 {
		t.Errorf("March expenses = %d, want 10850000", got.Report.Monthly[2].Expenses)
	}

	stored, err := svc.Get(created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if stored.Status != Completed {
		t.Errorf("stored status = %s, want completed", stored.Status)
	}
}

func TestProcessBalanceMismatchFailsTask(t *testing.T) {
	svc := newTestService()
	st := sampleStatement(t)
	st.ClosingBalance++
	created, err := svc.Create(st)
	if err != nil {
		t.Fatalf("balance mismatch must not block creation: %v", err)
	}

	got, err := svc.Process(created.ID)
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	if got.Status != Failed {
		t.Fatalf("status = %s, want failed", got.Status)
	}
	if got.FailReason == "" || got.Report != nil {
		t.Errorf("reason = %q, report = %v", got.FailReason, got.Report)
	}
}

func TestProcessTerminalTaskRejected(t *testing.T) {
	failing := sampleStatement(t)
	failing.ClosingBalance++

	tests := []struct {
		name string
		st   statement.Statement
	}{
		{"completed", sampleStatement(t)},
		{"failed", failing},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService()
			created, err := svc.Create(tt.st)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := svc.Process(created.ID); err != nil {
				t.Fatal(err)
			}

			_, err = svc.Process(created.ID)
			if !errors.Is(err, ErrInvalidTransition) {
				t.Fatalf("second Process() = %v, want ErrInvalidTransition", err)
			}
		})
	}
}

func TestProcessUnknownTask(t *testing.T) {
	svc := newTestService()
	if _, err := svc.Process("missing"); !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("Process(missing) = %v, want ErrTaskNotFound", err)
	}
}

func TestTasksAreIndependent(t *testing.T) {
	svc := newTestService()

	good, err := svc.Create(sampleStatement(t))
	if err != nil {
		t.Fatal(err)
	}
	other := statement.Statement{
		OpeningBalance: 1000,
		ClosingBalance: 400,
		Transactions: []statement.Transaction{
			{Date: time.Date(2026, 5, 3, 0, 0, 0, 0, time.UTC), Amount: -600, Currency: "KZT", Description: "Kinopark"},
		},
	}
	small, err := svc.Create(other)
	if err != nil {
		t.Fatal(err)
	}
	if good.ID == small.ID {
		t.Fatal("tasks share an ID")
	}

	if _, err := svc.Process(small.ID); err != nil {
		t.Fatal(err)
	}
	if got, _ := svc.Get(good.ID); got.Status != Pending || got.Report != nil {
		t.Errorf("processing one task touched another: %+v", got)
	}

	if _, err := svc.Process(good.ID); err != nil {
		t.Fatal(err)
	}
	gotSmall, _ := svc.Get(small.ID)
	gotGood, _ := svc.Get(good.ID)

	if len(gotSmall.Report.Monthly) != 1 || gotSmall.Report.Monthly[0].Month != "2026-05" {
		t.Errorf("small report mixed data: %+v", gotSmall.Report.Monthly)
	}
	if len(gotGood.Report.Monthly) != 3 {
		t.Errorf("good report mixed data: %+v", gotGood.Report.Monthly)
	}
	if len(gotSmall.Report.Subscriptions) != 0 || len(gotSmall.Report.Anomalies) != 0 {
		t.Errorf("small report picked up foreign results: %+v", gotSmall.Report)
	}
}

func TestMemoryStoreListSorted(t *testing.T) {
	s := NewMemoryStore()
	for _, id := range []string{"b", "c", "a"} {
		if err := s.Save(&Task{ID: id}); err != nil {
			t.Fatal(err)
		}
	}
	got := s.List()
	if len(got) != 3 || got[0].ID != "a" || got[1].ID != "b" || got[2].ID != "c" {
		t.Errorf("List() order = %v", got)
	}
}

func TestTaskTransitionKeepsStateOnError(t *testing.T) {
	tk := &Task{ID: "x", Status: Completed}
	if err := tk.fail("late"); !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("fail() = %v, want ErrInvalidTransition", err)
	}
	if tk.Status != Completed || tk.FailReason != "" {
		t.Errorf("task mutated on rejected transition: %+v", tk)
	}
}
