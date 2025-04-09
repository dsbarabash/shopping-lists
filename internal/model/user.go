package model

import (
	"github.com/google/uuid"
)

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	State    State  `json:"state"`
}

func NewUser(name, password string) *User {
	id := uuid.New()
	return &User{
		Id:       id.String(),
		Name:     name,
		Password: password,
		State:    1,
	}
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type RegistrationUserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
