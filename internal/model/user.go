package model

import (
	"github.com/google/uuid"
)

type User struct {
	id    string
	name  string
	state State
}

func NewUser(name string) *User {
	id := uuid.New()
	return &User{
		id:    id.String(),
		name:  name,
		state: 1,
	}
}
