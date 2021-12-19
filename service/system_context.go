package service

import (
	"github.com/kardianos/service"
	"github.com/rlaskowski/easymotion/config"
)

type SystemContext struct {
	service service.Service
}

func CreateSystemContext() (*SystemContext, error) {
	service, err := service.New(NewSystemService(), &service.Config{
		Name:             "EasyMotion",
		DisplayName:      "EasyMotion",
		Description:      "EasyMotion",
		WorkingDirectory: config.ProjectPath(),
		Option: service.KeyValue{
			"KeepAlive": true,
			"RunAtLoad": true,
		},
		Arguments: []string{
			"run",
		},
	})
	if err != nil {
		return nil, err
	}

	return &SystemContext{service}, nil
}

func (s *SystemContext) InstallService() error {
	return s.service.Install()
}

func (s *SystemContext) UninstallService() error {
	return s.service.Uninstall()
}

func (s *SystemContext) RestartService() error {
	return s.service.Restart()
}

func (s *SystemContext) RunService() error {
	return s.service.Run()
}
