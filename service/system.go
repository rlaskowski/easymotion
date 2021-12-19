package service

import "github.com/kardianos/service"

type SystemService struct {
	service *Service
}

func NewSystemService() *SystemService {
	return &SystemService{
		service: NewService(),
	}
}

func (s *SystemService) Start(srv service.Service) error {
	return s.service.Start()
}

func (s *SystemService) Stop(srv service.Service) error {
	return s.service.Stop()
}
