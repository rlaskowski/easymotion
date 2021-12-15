package service

import (
	"os"
	"os/signal"
	"syscall"

	"log"

	"github.com/rlaskowski/easymotion/config"
	"github.com/rlaskowski/easymotion/db"
)

type Service struct {
	httpServer    *HttpServer
	cameraService *CameraService
	sqlite        *db.Sqlite
	sigCh         chan os.Signal
}

func NewService() *Service {
	s := &Service{
		cameraService: NewCameraService(),
		sqlite:        db.NewSqlite(config.SqlitePath()),
		sigCh:         make(chan os.Signal, 1),
	}

	s.httpServer = NewHttpServer(s)

	return s
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

func (s *Service) CameraService() *CameraService {
	return s.cameraService
}

func (s *Service) Sqlite() *db.Sqlite {
	return s.sqlite
}
