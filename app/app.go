package app

import (
	"github.com/rlaskowski/easymotion"
	"github.com/rlaskowski/easymotion/auth"
	"github.com/rlaskowski/easymotion/service/dbservice"
	"github.com/rlaskowski/easymotion/service/opencvservice"
)

type App struct {
	immudb        *dbservice.ImmuDB
	opencvService *opencvservice.OpenCVService
}

func NewApp() *App {
	service, err := easymotion.GetService("service.opencv")
	if err != nil {
		return nil
	}

	openSrv := service.Intstance.(*opencvservice.OpenCVService)

	return &App{
		immudb:        dbservice.NewImmuDB(),
		opencvService: openSrv,
	}
}

func (a *App) CreateUser(name, email, password string) error {
	user := auth.NewUser(name, email, password)
	return a.immudb.CreateUser(user)
}

func (a *App) Users() ([]auth.User, error) {
	return a.immudb.Users()
}

func (a *App) OpenCVService() *opencvservice.OpenCVService {
	return a.opencvService
}
