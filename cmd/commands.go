package cmd

import (
	"log"
	"os"

	"github.com/rlaskowski/easymotion"
	"github.com/rlaskowski/easymotion/service/opencvservice"
	"github.com/rlaskowski/easymotion/service/queueservice"
	"github.com/rlaskowski/easymotion/service/storage"
	"github.com/rlaskowski/manage"
)

func RunCommand(service *easymotion.SystemService) {
	if len(os.Args) < 1 {
		log.Fatalf("use one of commnad: run | install | uninstall | start | stop | restart")
	}

	registerServices()

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
	default:
		log.Println("Please select option to run, for example: install | uninstall | restart | run")
		os.Exit(0)
	}
}

// Initializing services from service package
func registerServices() {
	manage.RegisterService(&opencvservice.OpenCVService{})
	manage.RegisterService(&queueservice.RabbitMQService{})
	manage.RegisterService(&storage.SqliteService{})
}
