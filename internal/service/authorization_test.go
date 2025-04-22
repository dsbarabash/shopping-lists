package service

import (
	"github.com/dsbarabash/shopping-lists/internal/model"
	"os"
	"testing"
)

func TestRegistration(t *testing.T) {
	type args struct {
		name     string
		password string
	}
	tests := []struct {
		name string
		args args
		want args
	}{
		{name: "ValidName", args: args{"Test", "password123"}, want: args{"Test", "password123"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Create("users.json")
			if err != nil {
				panic(err)
			}
			defer f.Close()
			defer func() {
				err := os.Remove("users.json")
				if err != nil {

				}
			}()

			if got := Registration(tt.args.name, tt.args.password); got.Name != tt.want.name || got.Password != tt.want.password {
				t.Errorf("Registration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	type args struct {
		user *model.User
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "Valid", args: args{user: &model.User{Name: "Test1", Password: "123", State: 1}}, want: "Test1", wantErr: false},
		{name: "InvalidStateArchived", args: args{user: &model.User{Name: "Test2", Password: "123", State: 0}}, want: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Create("users.json")
			if err != nil {
				panic(err)
			}
			defer f.Close()
			defer func() {
				err := os.Remove("users.json")
				if err != nil {

				}
			}()
			got, err := Login(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Login() got = %v, want %v", got, tt.want)
			}
		})
	}
}
