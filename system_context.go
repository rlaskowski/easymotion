package easymotion

import (
	"github.com/kardianos/service"
	"github.com/rlaskowski/manage"
)

type SystemContext struct {
	service *manage.ServiceInfo
	runner  Runner
}

func NewSystemContext() *SystemContext {
	return &SystemContext{
		service: &manage.ServiceInfo{},
		runner:  NewRunner(),
	}
}

func (s *SystemContext) Start(srv service.Service) error {
	if err := s.runner.RegisterServices(); err != nil {
		return err
	}

	if err := s.service.Start(); err != nil {
		return err
	}

	if err := s.runner.Run(); err != nil {
		return err
	}

	return nil
}

func (s *SystemContext) Stop(srv service.Service) error {
	if err := s.service.Stop(); err != nil {
		return err
	}

	if err := s.runner.Close(); err != nil {
		return err
	}

	return nil
}
