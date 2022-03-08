package dbservice

import (
	"fmt"
	"io/ioutil"

	"github.com/codenotary/immudb/embedded/sql"
	"github.com/codenotary/immudb/embedded/store"
	"github.com/rlaskowski/easymotion"
	"github.com/rlaskowski/easymotion/config"
	"github.com/rlaskowski/easymotion/embedded"
)

type ImmuDBService struct {
	name, path string
	immuStore  *store.ImmuStore
	engine     *sql.Engine
}

func (ImmuDBService) CreateService() *easymotion.ServiceInfo {
	return &easymotion.ServiceInfo{
		ID:        "service.database.immudb",
		Intstance: newImmuDBService(config.ProjectName(), config.ImmuDBPath()),
	}
}

func newImmuDBService(name, path string) *ImmuDBService {
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

	return i.initTable()
}

func (i *ImmuDBService) initTable() error {
	file, err := embedded.Files.Open("immudb/createtables.txt")
	if err != nil {
		return fmt.Errorf("couldn't open file with tables definition due to: %s", err.Error())
	}

	defer file.Close()

	script, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("bad definition of create tables: %s", err.Error())
	}

	_, _, err = i.engine.Exec(string(script), nil, nil)
	if err != nil {
		return fmt.Errorf("couldn't create initial tabels: %s", err.Error())
	}

	return nil
}
