package service

import (
	"log"

	"github.com/rlaskowski/easymotion/config"
)

type Service struct {
	httpServer     *HttpServer
	captureService *CaptureService
	immuDBService  *ImmuDBService
}

func NewService() *Service {
	s := &Service{
		captureService: NewCaptureService(),
		immuDBService:  NewImmuDBService(config.ImmuDBPath()),
	}

	s.httpServer = NewHttpServer(s)

	return s
}

func (s *Service) Start() error {
	var result error

	if result := s.captureService.Start(); result != nil {
		log.Printf("couldn't start capture service due to: %s", result.Error())
	}

	if result := s.immuDBService.Start(); result != nil {
		log.Printf("immudb dabase start problem: %s", result.Error())
	}

	if result := s.httpServer.Start(); result != nil {
		log.Printf("http server start problem: %s", result.Error())
	}

	if result == nil {
		log.Println("All services has been started")
	}

	return result
}

func (s *Service) Stop() error {
	log.Println("Stopping all services...")

	if err := s.captureService.Stop(); err != nil {
		log.Printf("couldn't stop campture service due to: %s", err.Error())
	}

	if err := s.immuDBService.Stop(); err != nil {
		log.Printf("immudb database stop problem: %s", err.Error())
	}

	if err := s.httpServer.Stop(); err != nil {
		log.Printf("http server stop problem: %s", err.Error())
	}

	return nil
}

func (s *Service) CaptureService() *CaptureService {
	return s.captureService
}
