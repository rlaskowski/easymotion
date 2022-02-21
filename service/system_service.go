package service

import (
	"github.com/kardianos/service"
)

type SystemService struct {
	service service.Service
}

func CreateSystemService(path string) (*SystemService, error) {
	service, err := service.New(NewSystemContext(), &service.Config{
		Name:             "easymotion",
		DisplayName:      "easymotion",
		Description:      "video capturing",
		WorkingDirectory: path,
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

func (s *SystemService) StartService() error {
	return s.service.Start()
}

func (s *SystemService) StopService() error {
	return s.service.Stop()
}
