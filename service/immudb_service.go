package service

import (
	"fmt"

	"github.com/codenotary/immudb/embedded/sql"
	"github.com/codenotary/immudb/embedded/store"
)

type ImmuDBService struct {
	name, path string
	immuStore  *store.ImmuStore
	engine     *sql.Engine
}

func NewImmuDBService(name, path string) *ImmuDBService {
	return &ImmuDBService{
		name: name,
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

	return i.initDB()
}

func (i *ImmuDBService) Stop() error {
	return i.immuStore.Close()
}

func (i *ImmuDBService) initDB() error {
	catalog, err := i.engine.Catalog(nil)
	if err != nil {
		return err
	}

	if !catalog.ExistDatabase(i.name) {
		_, _, err := i.engine.Exec(fmt.Sprintf("CREATE DATABASE %s", i.name), nil, nil)
		if err != nil {
			return err
		}
	}

	if err := i.engine.SetDefaultDatabase(i.name); err != nil {
		return err
	}

	return nil
}
