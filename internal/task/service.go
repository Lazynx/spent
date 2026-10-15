package task

import (
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
	return nil, errNotImplemented
}

func (s *service) Get(id string) (*Task, error) {
	return nil, errNotImplemented
}

func (s *service) Process(id string) (*Task, error) {
	return nil, errNotImplemented
}
