package logger

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
)

func TestParseLevel(t *testing.T) {
	tests := []struct {
		in   string
		want slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"INFO", slog.LevelInfo},
		{"warn", slog.LevelWarn},
		{"error", slog.LevelError},
		{"", slog.LevelInfo},
		{"bogus", slog.LevelInfo},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if got := parseLevel(tt.in); got != tt.want {
				t.Errorf("parseLevel(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestNewFormatAndLevel(t *testing.T) {
	tests := []struct {
		name     string
		level    string
		format   string
		wantSub  string
		wantNone bool
	}{
		{"json info", "info", "json", `"msg":"hello"`, false},
		{"text info", "info", "text", "msg=hello", false},
		{"debug filtered at warn", "warn", "json", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log := New(&buf, tt.level, tt.format)
			if tt.wantNone {
				log.Info("hello")
				if buf.Len() != 0 {
					t.Fatalf("expected no output, got %q", buf.String())
				}
				return
			}
			log.Info("hello")
			if !strings.Contains(buf.String(), tt.wantSub) {
				t.Errorf("output %q does not contain %q", buf.String(), tt.wantSub)
			}
		})
	}
}
