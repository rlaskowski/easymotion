package service

import (
	"github.com/codenotary/immudb/embedded/sql"
	"github.com/codenotary/immudb/embedded/store"
)

type ImmuDBService struct {
	path      string
	immuStore *store.ImmuStore
	engine    *sql.Engine
}

func NewImmuDBService(path string) *ImmuDBService {
	return &ImmuDBService{
		path: path,
	}
}

func (i *ImmuDBService) Start() error {
	immuStore, err := store.Open(i.path, store.DefaultOptions())
	if err != nil {
		return err
	}

	engine, err := sql.NewEngine(immuStore, sql.DefaultOptions())
	if err != nil {
		return err
	}

	i.immuStore = immuStore
	i.engine = engine

	return nil
}

func (i *ImmuDBService) Stop() error {
	return i.immuStore.Close()
}
