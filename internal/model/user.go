package model

type User struct {
	id    string
	name  string
	state State
}

func (u User) SetId(id string) {
	u.id = id
}

func (u User) GetId() string {
	return u.id
}

func (u User) SetName(name string) {
	u.name = name
}

func (u User) GetName() string {
	return u.name
}

func (u User) SetState(state State) {
	u.state = state
}

func (u User) GetState() string {
	return u.state.String()
}

func NewUser(id string, name string) *User {
	u := User{id: id, name: name}
	u.state = 1
	return &u
}
