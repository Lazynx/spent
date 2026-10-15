package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Lazynx/spent/internal/pkg/logger"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	logLevel := flag.String("log-level", "info", "debug | info | warn | error")
	logFormat := flag.String("log-format", "text", "text | json")
	flag.Parse()

	log := logger.New(os.Stderr, *logLevel, *logFormat)
	log.Info("spent started")

	return nil
}
