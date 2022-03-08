package dbservice

import (
	"fmt"

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
	user, err := i.User(user.Email)
	if err != nil {
		return err
	}

	if user != nil {
		return fmt.Errorf("user %s with email %s is already exist", user.Name, user.Email)
	}

	params := map[string]interface{}{
		"name":     user.Name,
		"email":    user.Email,
		"password": user.Password,
		"created":  user.Created,
	}

	_, _, err = i.engine.Exec(`INSERT INTO user(name,email,password,created) 
	                            VALUES (@name,@email,@password,@created)`, params, nil)

	return err

}

func (i *ImmuDB) User(email string) (*auth.User, error) {
	param := map[string]interface{}{
		"email": email,
	}

	query, err := i.engine.Query("SELECT u.name,u.email,u.password,u.created FROM user u WHERE u.email = @email", param, nil)
	if err != nil {
		return nil, err
	}

	user := &auth.User{}

	for {
		rows, err := query.Read()
		if err == sql.ErrNoMoreRows {
			break
		}

		fmt.Println(rows)

	}

	return user, nil

}

func (i *ImmuDB) Users() ([]*auth.User, error) {
	query, err := i.engine.Query("SELECT u.name,u.email,u.password,u.created FROM user u WHERE", nil, nil)
	if err != nil {
		return nil, err
	}

	users := make([]*auth.User, 0)

	for {
		rows, err := query.Read()
		if err == sql.ErrNoMoreRows {
			break
		}

		fmt.Println(rows)

	}

	return users, nil
}
