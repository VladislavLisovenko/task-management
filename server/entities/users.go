package entities

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (u User) GetID() int {
	return u.ID
}
