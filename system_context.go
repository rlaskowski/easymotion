package easymotion

import (
	"github.com/kardianos/service"
	"github.com/rlaskowski/manage"
)

type SystemContext struct {
	service *manage.ServiceInfo
}

func NewSystemContext() *SystemContext {
	return &SystemContext{
		service: &manage.ServiceInfo{},
	}
}

func (s *SystemContext) Start(srv service.Service) error {
	if err := s.service.Start(); err != nil {
		return err
	}

	return nil
}

func (s *SystemContext) Stop(srv service.Service) error {
	return s.service.Stop()
}
