package service

import "github.com/labstack/gommon/log"

type Service struct {
	httpServer *HttpServer
}

func NewService() *Service {
	return &Service{
		httpServer: NewHttpServer(),
	}
}

func (s *Service) Start() error {
	if err := s.httpServer.Stop(); err != nil {
		log.Infof("http server start problem: %s", err.Error())
	}
	return nil
}

func (s *Service) Stop() error {
	if err := s.httpServer.Stop(); err != nil {
		log.Infof("http server stop problem: %s", err.Error())
	}
	return nil
}
