package dbservice

import (
	"context"
	"fmt"
	"strings"

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
	params := map[string]interface{}{
		"email": email,
	}

	query, err := im.Query("SELECT u.name,u.email,u.created FROM user u WHERE u.email = @email", params)
	if err != nil {
		return auth.User{}, err
	}

	user := auth.User{}

	for query.Next() {

		if err := query.Scan(&user.Name, &user.Email, &user.Created); err != nil {
			return auth.User{}, err
		}

	}

	return user, nil
}

func (im *ImmuDB) Users() ([]auth.User, error) {
	query, err := im.Query("SELECT u.name,u.email,u.created FROM user u", nil)
	if err != nil {
		return nil, err
	}

	users := make([]auth.User, 0)

	for query.Next() {

		user := auth.User{}

		if err := query.Scan(&user.Name, &user.Email, &user.Created); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (im *ImmuDB) Query(sql string, params map[string]interface{}) (*ImmuDBRows, error) {
	return im.QueryContext(context.Background(), sql, params)
}

func (im *ImmuDB) QueryContext(ctx context.Context, sql string, params map[string]interface{}) (*ImmuDBRows, error) {
	rowr, err := im.engine.Query(sql, params, nil)

	if err != nil {
		return nil, err
	}

	rows := &ImmuDBRows{
		rowr: rowr,
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	return rows, nil

}
