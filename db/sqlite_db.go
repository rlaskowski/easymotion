package db

import (
	"database/sql"
	"log"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	name, path string
	database   *sql.DB
}

func NewSqlite(path string) *Sqlite {
	return &Sqlite{
		name: filepath.Base(path),
		path: path,
	}
}

func (s *Sqlite) Start() error {
	log.Println("Starting sqlite database")

	database, err := sql.Open("sqlite3", s.path)
	if err != nil {
		return err
	}

	s.database = database

	return nil
}

func (s *Sqlite) Stop() error {
	log.Println("Stopping sqlite database")
	return s.database.Close()
}
