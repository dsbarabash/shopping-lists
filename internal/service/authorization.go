package service

import (
	"errors"
	"github.com/dsbarabash/shopping-lists/internal/config"
	"github.com/dsbarabash/shopping-lists/internal/model"
	"github.com/dsbarabash/shopping-lists/internal/repository"
	"github.com/golang-jwt/jwt/v5"
)

func Registration(name, password string) *model.User {
	var user *model.User
	user = model.NewUser(name, password)
	repository.UserList.SaveToFile(user)
	return user
}

func Login(user *model.User) (string, error) {
	repository.UserList.Mu.Lock()
	defer repository.UserList.Mu.Unlock()
	if user.State != 1 {
		return "", errors.New("USER NOT ACTIVE")
	}
	for _, u := range repository.UserList.Store {
		if u.Id == user.Id {
			if u.Password == user.Password && u.Name == user.Name {
				secretKey := []byte(config.Cfg.Secret)
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"id":       user.Id,
					"name":     user.Name,
					"password": user.Password,
					"state":    1,
				})
				tokenString, err := token.SignedString(secretKey)
				if err != nil {
					return "", err
				}
				return tokenString, nil
			} else {
				return "", errors.New("LOGIN OR PASSWORD IS INCORRECT")
			}

		}
	}
	return "", errors.New("NOT FOUND")

}
