package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rlaskowski/easymotion/config"
	"github.com/rlaskowski/easymotion/service"
)

func init() {
	service.RegisterService(&service.HttpServer{})
	service.RegisterService(&service.CaptureService{})
	service.RegisterService(&service.ImmuDBService{})
}

func main() {
	systemService, err := service.CreateSystemService(config.ProjectPath())
	if err != nil {
		log.Fatalf("Unexpected error: %s", err.Error())
	}

	runCommand(systemService)
}

func runCommand(service *service.SystemService) {
	if len(os.Args) < 2 {
		fmt.Println("Please select option to run, for example: install | uninstall | restart | run")
	}

	switch os.Args[1] {
	case "run":
		if err := service.RunService(); err != nil {
			log.Println(err)
		}
	case "install":
		if err := service.StopService(); err != nil {
			log.Println(err)
		}

		if err := service.UninstallService(); err != nil {
			log.Println(err)
		}

		if err := service.InstallService(); err != nil {
			log.Println(err)
		}

		if err := service.StartService(); err != nil {
			log.Println(err)
		}
	case "uninstall":
		if err := service.StopService(); err != nil {
			log.Println(err)
		}

		if err := service.UninstallService(); err != nil {
			log.Println(err)
		}
	case "start":
		if err := service.StartService(); err != nil {
			log.Println(err)
		}
	case "stop":
		if err := service.StopService(); err != nil {
			log.Println(err)
		}
	case "restart":
		if err := service.RestartService(); err != nil {
			log.Println(err)
		}
	}
}
