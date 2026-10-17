package task

import "slices"

type Status int

const (
	Pending Status = iota
	Processing
	Completed
	Failed
)

var statusNames = [...]string{"pending", "processing", "completed", "failed"}

var transitions = map[Status][]Status{
	Pending:    {Processing},
	Processing: {Completed, Failed},
}

func (s Status) String() string {
	if s < 0 || int(s) >= len(statusNames) {
		return "unknown"
	}
	return statusNames[s]
}

func (s Status) CanTransitionTo(next Status) bool {
	return slices.Contains(transitions[s], next)
}
