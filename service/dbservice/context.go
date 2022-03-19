package dbservice

import "github.com/rlaskowski/easymotion/auth"

type Context interface {
	CreateUser(user *auth.User) error
	UserByEmail(email string) (auth.User, error)
	Users() ([]auth.User, error)
}
