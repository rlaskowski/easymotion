package service

import (
	"log"

	"github.com/rlaskowski/easymotion/config"
	"github.com/rlaskowski/easymotion/db"
)

type Service struct {
	httpServer     *HttpServer
	captureService *CaptureService
	sqlite         *db.Sqlite
}

func NewService() *Service {
	s := &Service{
		captureService: NewCaptureService(),
		sqlite:         db.NewSqlite(config.SqlitePath()),
	}

	s.httpServer = NewHttpServer(s)

	return s
}

func (s *Service) Start() error {
	var result error

	if result := s.captureService.Start(); result != nil {
		log.Printf("couldn't start capture service due to: %s", result.Error())
	}

	if result := s.sqlite.Start(); result != nil {
		log.Printf("sqlite dabase start problem: %s", result.Error())
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

	if err := s.sqlite.Stop(); err != nil {
		log.Printf("sqlite database stop problem: %s", err.Error())
	}

	if err := s.httpServer.Stop(); err != nil {
		log.Printf("http server stop problem: %s", err.Error())
	}

	return nil
}

func (s *Service) CaptureService() *CaptureService {
	return s.captureService
}

func (s *Service) Sqlite() *db.Sqlite {
	return s.sqlite
}
