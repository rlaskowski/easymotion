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
	var result []error

	if err := s.sqlite.Start(); err != nil {
		result = append(result, err)
		log.Printf("sqlite dabase start problem: %s", err.Error())
	}

	if err := s.httpServer.Start(); err != nil {
		result = append(result, err)
		log.Printf("http server start problem: %s", err.Error())
	}

	if result != nil {
		os.Exit(0)
	} else {
		log.Println("All services has been started")
	}

	signal.Notify(s.sigCh, syscall.SIGINT|syscall.SIGTERM|syscall.SIGKILL)

	<-s.sigCh

	return nil
}

func (s *Service) Stop() error {
	log.Println("Stopping all services...")

	if err := s.sqlite.Stop(); err != nil {
		log.Printf("sqlite database stop problem: %s", err.Error())
	}

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
