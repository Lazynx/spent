package task

import (
	"fmt"

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
	if !t.Status.CanTransitionTo(next) {
		return fmt.Errorf("task %s: %s -> %s: %w", t.ID, t.Status, next, ErrInvalidTransition)
	}
	t.Status = next
	return nil
}

func (t *Task) start() error {
	return t.transition(Processing)
}

func (t *Task) complete(r analyzer.Report) error {
	if err := t.transition(Completed); err != nil {
		return err
	}
	t.Report = &r
	return nil
}

func (t *Task) fail(reason string) error {
	if err := t.transition(Failed); err != nil {
		return err
	}
	t.FailReason = reason
	return nil
}
