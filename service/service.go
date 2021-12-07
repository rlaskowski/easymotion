package service

import (
	"os"
	"os/signal"
	"syscall"

	"log"
)

type Service struct {
	httpServer *HttpServer
	sigCh      chan os.Signal
}

func NewService() *Service {
	return &Service{
		httpServer: NewHttpServer(),
		sigCh:      make(chan os.Signal, 1),
	}
}

func (s *Service) Start() error {
	log.Println("Starting all services...")

	if err := s.httpServer.Start(); err != nil {
		log.Printf("http server start problem: %s", err.Error())
	}

	signal.Notify(s.sigCh, syscall.SIGINT|syscall.SIGTERM|syscall.SIGKILL)

	<-s.sigCh

	return nil
}

func (s *Service) Stop() error {
	log.Println("Stopping all services...")

	if err := s.httpServer.Stop(); err != nil {
		log.Printf("http server stop problem: %s", err.Error())
	}
	return nil
}
