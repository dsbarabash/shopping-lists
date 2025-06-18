package model

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	State    State  `json:"state"`
}

func NewUser(dto *CreateUserDTO) *User {
	return &User{
		Id:       dto.Id,
		Name:     dto.Name,
		Password: dto.Password,
		State:    dto.State,
	}
}

type CreateUserDTO struct {
	Id       string `json:"id,omitempty" bson:"id,omitempty"`
	Name     string `json:"name,omitempty" bson:"name,omitempty"`
	Password string `json:"password,omitempty" bson:"password,omitempty"`
	State    State  `json:"state,omitempty" bson:"state,omitempty"`
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type RegistrationUserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	Success string `json:"name"`
	Token   string `json:"token"`
}
