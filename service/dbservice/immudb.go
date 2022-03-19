package dbservice

import (
	"fmt"
	"strings"
	"time"

	"github.com/codenotary/immudb/embedded/sql"
	"github.com/rlaskowski/easymotion"
	"github.com/rlaskowski/easymotion/auth"
)

type ImmuDB struct {
	engine *sql.Engine
}

func NewImmuDB() *ImmuDB {
	service, err := easymotion.GetService("service.database.immudb")
	if err != nil {
		return nil
	}

	db := service.Intstance.(*ImmuDBService)

	return &ImmuDB{
		engine: db.Engine(),
	}
}

func (im *ImmuDB) CreateUser(user *auth.User) error {
	u, err := im.UserByEmail(user.Email)
	if err != nil {
		return err
	}

	if strings.Compare(u.Email, user.Email) == 0 {
		return fmt.Errorf("user %s with email %s is already exist", user.Name, user.Email)
	}

	params := map[string]interface{}{
		"name":     user.Name,
		"email":    user.Email,
		"password": user.Password,
		"created":  user.Created,
	}

	_, _, err = im.engine.Exec(`INSERT INTO user(name,email,password,created) 
	                            VALUES (@name,@email,@password,@created)`, params, nil)

	return err

}

func (im *ImmuDB) UserByEmail(email string) (auth.User, error) {
	param := map[string]interface{}{
		"email": email,
	}

	query, err := im.engine.Query("SELECT u.name,u.email,u.created FROM user u WHERE u.email = @email", param, nil)
	if err != nil {
		return auth.User{}, err
	}

	user := auth.User{}

	for {
		row, err := query.Read()
		if err == sql.ErrNoMoreRows {
			break
		}

		if err := im.scan(row, &user.Name, &user.Email, &user.Created); err != nil {
			return auth.User{}, err
		}

	}

	return user, nil
}

func (im *ImmuDB) Users() ([]auth.User, error) {
	query, err := im.engine.Query("SELECT u.name,u.email,u.created FROM user u", nil, nil)
	if err != nil {
		return nil, err
	}

	users := make([]auth.User, 0)

	for {
		row, err := query.Read()
		if err == sql.ErrNoMoreRows {
			break
		}

		user := auth.User{}

		if err := im.scan(row, &user.Name, &user.Email, &user.Created); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (im *ImmuDB) scan(row *sql.Row, params ...interface{}) error {
	var cols []interface{}

	for _, r := range row.Values {
		cols = append(cols, r.Value())
	}

	if len(cols) != len(params) {
		return fmt.Errorf("different number of columns in row, expected %d got %d", len(params), len(cols))
	}

	for i, v := range cols {
		if err := im.parseType(params[i], v); err != nil {
			return err
		}
	}

	return nil
}

func (im *ImmuDB) parseType(dst, src interface{}) error {
	switch s := src.(type) {
	case string:
		switch d := dst.(type) {
		case *string:
			*d = s
			return nil
		}
	case time.Time:
		switch d := dst.(type) {
		case *time.Time:
			*d = s
			return nil
		case *string:
			*d = s.Format(time.RFC3339Nano)
			return nil
		case *[]byte:
			*d = []byte(s.Format(time.RFC3339Nano))
			return nil
		}
	}

	return nil
}
