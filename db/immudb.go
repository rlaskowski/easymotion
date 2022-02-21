package db

import (
	"errors"

	"github.com/codenotary/immudb/embedded/sql"
	"github.com/rlaskowski/easymotion/auth"
)

type ImmuDB struct {
	engine *sql.Engine
}

func NewImmuDB(engine *sql.Engine) *ImmuDB {
	return &ImmuDB{engine}
}

func (i *ImmuDB) CreateUser(user *auth.User) error {
	return errors.New("method has been not yet implemented")
}

func (i *ImmuDB) User(email string) (*auth.User, error) {
	return nil, errors.New("method has been not yet implemented")
}

func (i *ImmuDB) Users() ([]*auth.User, error) {
	return nil, errors.New("method has been not yet implemented")
}
