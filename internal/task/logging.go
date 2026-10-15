package task

import (
	"log/slog"

	"github.com/Lazynx/spent/internal/statement"
)

type loggingService struct {
	next Service
	log  *slog.Logger
}

func WithLogging(next Service, log *slog.Logger) Service {
	return &loggingService{next: next, log: log}
}

func (s *loggingService) Create(st statement.Statement) (*Task, error) {
	return s.next.Create(st)
}

func (s *loggingService) Get(id string) (*Task, error) {
	return s.next.Get(id)
}

func (s *loggingService) Process(id string) (*Task, error) {
	return s.next.Process(id)
}
