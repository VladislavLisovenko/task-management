package entities

import "time"

type TaskListFilter struct {
	User               User
	Description        string
	ExpirationDateFrom time.Time
	ExpirationDateTo   time.Time
	Done               bool
	ConsiderDone       bool
	FilterIsSet        bool
}

func (tlf *TaskListFilter) SetDescription(v string) {
	tlf.Description = v
	tlf.FilterIsSet = true
}

func (tlf *TaskListFilter) SetExpirationDateFrom(v time.Time) {
	tlf.ExpirationDateFrom = v
	tlf.FilterIsSet = true
}

func (tlf *TaskListFilter) SetExpirationDateTo(v time.Time) {
	tlf.ExpirationDateTo = v
	tlf.FilterIsSet = true
}

func (tlf *TaskListFilter) SetDone(v bool) {
	tlf.Done = v
	tlf.ConsiderDone = true
	tlf.FilterIsSet = true
}
