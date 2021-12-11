package db

import "errors"

type SqliteDB struct {
}

func (s *SqliteDB) Start() error {
	return errors.New("not yet implemented")
}

func (s *SqliteDB) Stop() error {
	return errors.New("not yet implemented")
}
