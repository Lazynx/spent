package task

type Status int

const (
	Pending Status = iota
	Processing
	Completed
	Failed
)

func (s Status) String() string {
	return ""
}

func (s Status) CanTransitionTo(next Status) bool {
	return false
}
