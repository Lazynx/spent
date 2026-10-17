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
	t, err := s.next.Create(st)
	if err != nil {
		s.log.Warn("create task failed", "transactions", len(st.Transactions), "error", err)
		return nil, err
	}
	s.log.Info("task created", "id", t.ID, "transactions", len(st.Transactions))
	return t, nil
}

func (s *loggingService) Get(id string) (*Task, error) {
	t, err := s.next.Get(id)
	if err != nil {
		s.log.Warn("get task failed", "id", id, "error", err)
		return nil, err
	}
	s.log.Debug("task fetched", "id", id, "status", t.Status.String())
	return t, nil
}

func (s *loggingService) Process(id string) (*Task, error) {
	t, err := s.next.Process(id)
	if err != nil {
		s.log.Warn("process task failed", "id", id, "error", err)
		return nil, err
	}
	if t.Status == Failed {
		s.log.Warn("task processing failed", "id", id, "reason", t.FailReason)
		return t, nil
	}
	s.log.Info("task processed", "id", id, "status", t.Status.String())
	return t, nil
}
