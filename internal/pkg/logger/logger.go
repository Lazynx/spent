package logger

import (
	"io"
	"log/slog"
	"strings"
)

func New(w io.Writer, level, format string) *slog.Logger {
	opts := &slog.HandlerOptions{Level: parseLevel(level)}

	var h slog.Handler
	if format == "json" {
		h = slog.NewJSONHandler(w, opts)
	} else {
		h = slog.NewTextHandler(w, opts)
	}
	return slog.New(h)
}

func Discard() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func parseLevel(s string) slog.Level {
	var lvl slog.Level
	if err := lvl.UnmarshalText([]byte(strings.ToLower(s))); err != nil {
		return slog.LevelInfo
	}
	return lvl
}
