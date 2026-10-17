package task

import (
	"cmp"
	"fmt"
	"slices"
)

type Store interface {
	Save(t *Task) error
	Get(id string) (*Task, error)
	List() []*Task
}

type MemoryStore struct {
	tasks map[string]*Task
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{tasks: make(map[string]*Task)}
}

func (s *MemoryStore) Save(t *Task) error {
	s.tasks[t.ID] = t
	return nil
}

func (s *MemoryStore) Get(id string) (*Task, error) {
	t, ok := s.tasks[id]
	if !ok {
		return nil, fmt.Errorf("task %q: %w", id, ErrTaskNotFound)
	}
	return t, nil
}

func (s *MemoryStore) List() []*Task {
	out := make([]*Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		out = append(out, t)
	}
	slices.SortFunc(out, func(a, b *Task) int { return cmp.Compare(a.ID, b.ID) })
	return out
}
