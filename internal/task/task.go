package task

import (
	"github.com/Lazynx/spent/internal/analyzer"
	"github.com/Lazynx/spent/internal/statement"
)

type Task struct {
	ID         string
	Statement  statement.Statement
	Status     Status
	Report     *analyzer.Report
	FailReason string
}

func (t *Task) transition(next Status) error {
	return errNotImplemented
}

func (t *Task) start() error {
	return errNotImplemented
}

func (t *Task) complete(r analyzer.Report) error {
	return errNotImplemented
}

func (t *Task) fail(reason string) error {
	return errNotImplemented
}
