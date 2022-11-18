package storage

import (
	"database/sql"
	"log"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rlaskowski/easymotion/config"
	"github.com/rlaskowski/manage"
)

type SqliteService struct {
	name, path string
	database   *sql.DB
}

func (SqliteService) CreateService() *manage.ServiceInfo {
	return &manage.ServiceInfo{
		ID:        "service.sqlite",
		Priority:  3,
		Intstance: newSqlite(),
	}
}

func newSqlite() *SqliteService {
	return &SqliteService{
		name: "easymotion",
		path: filepath.Join(config.WorkingDirectory(), "sqlite"),
	}
}

func (s *SqliteService) Start() error {
	log.Println("starting sqlite service")

	database, err := sql.Open("sqlite3", s.path)
	if err != nil {
		return err
	}

	s.database = database

	return nil
}

func (s *SqliteService) Stop() error {
	log.Println("stopping sqlite service")

	return s.database.Close()
}
