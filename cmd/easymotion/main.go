package main

import (
	"log"
	"os"

	"github.com/rlaskowski/easymotion"
	"github.com/rlaskowski/easymotion/cmd"
	"github.com/rlaskowski/easymotion/config"
	"github.com/rlaskowski/easymotion/service/dbservice"
	"github.com/rlaskowski/easymotion/service/httpservice"
	"github.com/rlaskowski/easymotion/service/opencvservice"
)

// List of services that need to be run from the Service type
func init() {
	easymotion.RegisterService(&httpservice.HttpServer{})
	easymotion.RegisterService(&opencvservice.OpenCVService{})
	easymotion.RegisterService(&dbservice.ImmuDBService{})
}

func main() {
	systemService, err := easymotion.CreateSystemService(config.ProjectPath())
	if err != nil {
		log.Fatalf("Unexpected error: %s", err.Error())
		os.Exit(1)
	}

	cmd.RunCommand(systemService)
}
