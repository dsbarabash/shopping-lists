package model

import (
	"github.com/google/uuid"
)

type User struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	State State  `json:"state"`
}

func NewUser(name string) *User {
	id := uuid.New()
	return &User{
		Id:    id.String(),
		Name:  name,
		State: 1,
	}
}
