package task

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
	return errNotImplemented
}

func (s *MemoryStore) Get(id string) (*Task, error) {
	return nil, errNotImplemented
}

func (s *MemoryStore) List() []*Task {
	return nil
}
