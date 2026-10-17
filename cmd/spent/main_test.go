package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestFormatMoney(t *testing.T) {
	tests := []struct {
		in   int64
		want string
	}{
		{0, "0.00"},
		{5, "0.05"},
		{250000, "2500.00"},
		{-1550, "-15.50"},
	}
	for _, tt := range tests {
		if got := formatMoney(tt.in); got != tt.want {
			t.Errorf("formatMoney(%d) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestRunSample(t *testing.T) {
	var out, errOut bytes.Buffer
	err := run([]string{"-log-level", "error", "../../testdata/sample.csv"}, &out, &errOut)
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}

	got := strings.Join(strings.Fields(out.String()), " ")
	for _, want := range []string{
		"2026-03 income 500000.00 expenses 108500.00 net 391500.00",
		"Food 156000.00",
		"Entertainment 7500.00",
		"magnum 139000.00",
		"netflix.com 2500.00 x3, last 2026-03-13",
		"2026-03-15 Magnum 90000.00 (Food median 15500.00)",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("output missing %q\n%s", want, got)
		}
	}
}

func TestRunUsageError(t *testing.T) {
	var out, errOut bytes.Buffer
	if err := run(nil, &out, &errOut); err == nil {
		t.Fatal("expected usage error")
	}
}

func TestRunMissingFile(t *testing.T) {
	var out, errOut bytes.Buffer
	if err := run([]string{"nope.csv"}, &out, &errOut); err == nil {
		t.Fatal("expected error for missing file")
	}
}
