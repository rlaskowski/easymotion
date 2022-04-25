package dbservice

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/codenotary/immudb/embedded/sql"
	"github.com/rlaskowski/easymotion"
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

func (im *ImmuDB) CreateUser(user *User) error {
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

func (im *ImmuDB) UserByEmail(email string) (User, error) {
	params := map[string]interface{}{
		"email": email,
	}

	query, err := im.Query("SELECT u.name,u.email,u.created FROM user u WHERE u.email = @email", params)
	if err != nil {
		return User{}, err
	}

	user := User{}

	for query.Next() {

		if err := query.Scan(&user.Name, &user.Email, &user.Created); err != nil {
			return User{}, err
		}

	}

	return user, nil
}

func (im *ImmuDB) Users() ([]User, error) {
	query, err := im.Query("SELECT u.name,u.email,u.created FROM user u", nil)
	if err != nil {
		return nil, err
	}

	users := make([]User, 0)

	for query.Next() {

		user := User{}

		if err := query.Scan(&user.Name, &user.Email, &user.Created); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (im *ImmuDB) CreateCamOption(options *CameraOptions) error {
	cop, err := im.CameraOption(options.CameraID)
	if err != nil {
		return err
	}

	if reflect.DeepEqual(options.CameraID, cop.CameraID) {
		return fmt.Errorf("camera name %s is already exist", options.Name)
	}

	params := map[string]interface{}{
		"name":     options.Name,
		"cameraID": options.CameraID,
		"autorec":  options.Autorec,
		"timeline": options.Timeline,
	}

	_, _, err = im.engine.Exec(`INSERT INTO camera_options(name,camera_id,auto_recording,timeline) 
	                            VALUES (@name,@cameraID,@autorec,@timeline)`, params, nil)

	return err
}

func (im *ImmuDB) CameraOptions() ([]CameraOptions, error) {
	sql := `SELECT c.id, c.name, c.camera_id, c.auto_recording, c.timeline 
	        FROM camera_options`

	options, err := im.cameraOptions(sql, nil)

	if err != nil {
		return nil, err
	}

	return options, nil
}

func (im *ImmuDB) CameraOption(camID int) (CameraOptions, error) {
	params := map[string]interface{}{
		"camID": camID,
	}

	sql := `SELECT c.id, c.name, c.camera_id, c.auto_recording, c.timeline 
	        FROM camera_options c
	        WHERE c.camera_id = @camID`

	options, err := im.cameraOptions(sql, params)

	if err != nil {
		return CameraOptions{}, err
	}

	lopt := len(options)

	if lopt > 0 {
		return CameraOptions{}, errors.New("unexpected option")
	}

	if cap(options) == 0 {
		return CameraOptions{CameraID: -1}, nil
	}

	return options[0], nil
}

func (im *ImmuDB) cameraOptions(sql string, params map[string]interface{}) ([]CameraOptions, error) {
	query, err := im.Query(sql, params)

	if err != nil {
		return nil, err
	}

	options := make([]CameraOptions, 0)

	for query.Next() {
		cop := CameraOptions{}

		if err := query.Scan(&cop.ID, &cop.Name, &cop.CameraID, &cop.Autorec, &cop.Timeline); err != nil {
			return nil, err
		}

		options = append(options, cop)
	}

	return options, nil
}

// Selecting data by params
func (im *ImmuDB) Query(sql string, params map[string]interface{}) (*ImmuDBRows, error) {
	return im.QueryContext(context.Background(), sql, params)
}

// Selecting data by params with context
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
