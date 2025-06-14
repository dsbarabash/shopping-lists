package repository

import (
	"github.com/google/uuid"
)

//go:generate mockgen -source=id.go -destination=id_mock.go -package=repository

type Ider interface {
	Id() (string, error)
}

type idImpl struct{}

func NewIder() Ider {
	return &idImpl{}
}

func (i idImpl) Id() (string, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return uuid.String(), nil
}
