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
