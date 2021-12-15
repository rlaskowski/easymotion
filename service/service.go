package service

import (
	"os"
	"os/signal"
	"syscall"

	"log"

	"github.com/rlaskowski/easymotion/db"
)

type Service struct {
	httpServer    *HttpServer
	cameraService *CameraService
	sqliteDB      *db.SqliteDB
	sigCh         chan os.Signal
}

func NewService() *Service {
	s := &Service{
		cameraService: NewCameraService(),
		sqliteDB:      db.NewSqliteDB("./easymotion.db"),
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

func (s *Service) SqliteDB() *db.SqliteDB {
	return s.sqliteDB
}
