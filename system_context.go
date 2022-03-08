package easymotion

import (
	"os"

	"github.com/kardianos/service"
)

type SystemContext struct {
	service *Service
}

func NewSystemContext() *SystemContext {
	return &SystemContext{
		service: NewService(),
	}
}

func (s *SystemContext) Start(srv service.Service) error {
	go func() {
		if err := s.service.Start(); err != nil {
			os.Exit(1)
		}

	}()
	return nil
}

func (s *SystemContext) Stop(srv service.Service) error {
	return s.service.Stop()
}
