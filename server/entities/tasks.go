package entities

import "time"

type Task struct {
	ID             int       `json:"id"`
	Description    string    `json:"description"`
	ExpirationDate time.Time `json:"expirationDate"`
	CreationDate   time.Time `json:"creationDate"`
	User           User      `json:"user"`
	Done           bool      `json:"done"`
}

func (o Task) GetID() int {
	return o.ID
}
