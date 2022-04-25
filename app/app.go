package app

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/rlaskowski/easymotion"
	"github.com/rlaskowski/easymotion/service/dbservice"
	"github.com/rlaskowski/easymotion/service/opencvservice"
)

type App struct {
	dbctx         dbservice.Context
	opencvService *opencvservice.OpenCVService
}

func NewApp() *App {
	service, err := easymotion.GetService("service.opencv")
	if err != nil {
		panic(err)
	}

	openSrv := service.Intstance.(*opencvservice.OpenCVService)

	return &App{
		dbctx:         dbservice.NewImmuDB(),
		opencvService: openSrv,
	}
}

func (a *App) camera(id int) (*opencvservice.Camera, error) {
	options, err := a.dbctx.CameraOption(id)

	if err != nil {
		return nil, fmt.Errorf("no options for camera %d", id)
	}

	opt := opencvservice.CameraOptions{
		Autorec:  options.Autorec,
		Timeline: options.Timeline,
	}

	camera, err := a.opencvService.GetOrCreate(id, opt)

	if err != nil {
		return nil, fmt.Errorf("camera %d not found", id)
	}

	return camera, nil
}

func (a *App) ReadBytes(id int) ([]byte, error) {
	camera, err := a.camera(id)

	if err != nil {
		return nil, err
	}

	buff := bytes.Buffer{}

	if _, err := io.Copy(&buff, camera); err != nil {
		return nil, errors.New("nothing to read")
	}

	return buff.Bytes(), nil
}

func (a *App) StartRecord(id int) error {
	camera, err := a.camera(id)

	if err != nil {
		return err
	}

	return camera.StopRecord()
}

func (a *App) StopRecord(id int) error {
	camera, err := a.camera(id)

	if err != nil {
		return err
	}

	return camera.StopRecord()
}

func (a *App) Users() ([]dbservice.User, error) {
	return a.dbctx.Users()
}

func (a *App) CreateUser(name, email, password string) error {
	user := dbservice.NewUser(name, email, password)
	return a.dbctx.CreateUser(user)
}
