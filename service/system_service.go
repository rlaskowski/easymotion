package service

import (
	"github.com/kardianos/service"
	"github.com/rlaskowski/easymotion/config"
)

type SystemService struct {
	service service.Service
}

func CreateSystemService() (*SystemService, error) {
	service, err := service.New(NewSystemContext(), &service.Config{
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

	return &SystemService{service}, nil
}

func (s *SystemService) InstallService() error {
	return s.service.Install()
}

func (s *SystemService) UninstallService() error {
	return s.service.Uninstall()
}

func (s *SystemService) RestartService() error {
	return s.service.Restart()
}

func (s *SystemService) RunService() error {
	return s.service.Run()
}

func (s *SystemService) StatusService() (service.Status, error) {
	return s.service.Status()
}
