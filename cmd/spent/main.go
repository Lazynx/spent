package main

import (
	"crypto/rand"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Lazynx/spent/internal/analyzer"
	"github.com/Lazynx/spent/internal/pkg/logger"
	"github.com/Lazynx/spent/internal/statement"
	"github.com/Lazynx/spent/internal/task"
)

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("spent", flag.ContinueOnError)
	fs.SetOutput(stderr)
	logLevel := fs.String("log-level", "info", "debug | info | warn | error")
	logFormat := fs.String("log-format", "text", "text | json")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return errors.New("usage: spent [flags] <statement.csv>")
	}

	log := logger.New(stderr, *logLevel, *logFormat)

	st, err := parseFile(statement.CSVParser{}, fs.Arg(0))
	if err != nil {
		return err
	}

	categorizer := analyzer.NewRuleCategorizer()
	analyze := func(s statement.Statement) (analyzer.Report, error) {
		return analyzer.Analyze(s, categorizer)
	}
	svc := task.WithLogging(
		task.NewService(task.NewMemoryStore(), analyze, time.Now, rand.Text),
		log,
	)

	created, err := svc.Create(st)
	if err != nil {
		return err
	}
	processed, err := svc.Process(created.ID)
	if err != nil {
		return err
	}
	if processed.Status == task.Failed {
		return fmt.Errorf("task %s failed: %s", processed.ID, processed.FailReason)
	}

	return printReport(stdout, *processed.Report)
}

func parseFile(p statement.Parser, path string) (statement.Statement, error) {
	f, err := os.Open(path)
	if err != nil {
		return statement.Statement{}, err
	}
	defer f.Close()

	st, err := p.Parse(f)
	if err != nil {
		return statement.Statement{}, fmt.Errorf("parse %s: %w", path, err)
	}
	return st, nil
}
