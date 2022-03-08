package auth

import (
	"time"
)

type User struct {
	ID       int16     `json:"-"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
	Created  time.Time `json:"created"`
}

func NewUser(name, email, password string) *User {
	return &User{
		Name:     name,
		Email:    email,
		Password: password,
		Created:  time.Now(),
	}
}
