package task

import (
	"fmt"
	"slices"
	"time"

	"github.com/Lazynx/spent/internal/analyzer"
	"github.com/Lazynx/spent/internal/statement"
)

type Service interface {
	Create(st statement.Statement) (*Task, error)
	Get(id string) (*Task, error)
	Process(id string) (*Task, error)
}

type Analyzer func(statement.Statement) (analyzer.Report, error)

type service struct {
	store   Store
	analyze Analyzer
	now     func() time.Time
	newID   func() string
}

func NewService(store Store, analyze Analyzer, now func() time.Time, newID func() string) Service {
	return &service{store: store, analyze: analyze, now: now, newID: newID}
}

func (s *service) Create(st statement.Statement) (*Task, error) {
	if err := st.Validate(s.now()); err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}

	st.Transactions = slices.Clone(st.Transactions)
	t := &Task{
		ID:        s.newID(),
		Statement: st,
		Status:    Pending,
	}
	if err := s.store.Save(t); err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}
	return t, nil
}

func (s *service) Get(id string) (*Task, error) {
	return s.store.Get(id)
}

func (s *service) Process(id string) (*Task, error) {
	t, err := s.store.Get(id)
	if err != nil {
		return nil, err
	}

	if err := t.start(); err != nil {
		return nil, err
	}
	if err := s.store.Save(t); err != nil {
		return nil, err
	}

	report, err := s.analyze(t.Statement)
	if err != nil {
		if ferr := t.fail(err.Error()); ferr != nil {
			return nil, ferr
		}
	} else if cerr := t.complete(report); cerr != nil {
		return nil, cerr
	}

	if err := s.store.Save(t); err != nil {
		return nil, err
	}
	return t, nil
}
