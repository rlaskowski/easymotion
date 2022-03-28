package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/rlaskowski/easymotion"
	"github.com/rlaskowski/easymotion/config"
)

var VideosPath string

func flags() {
	flag.StringVar(&VideosPath, "f", config.ProjectPath(), "Path where will be store video files")

	flag.Parse()

}

func RunCommand(service *easymotion.SystemService) {
	if len(os.Args) < 2 {
		fmt.Println("Please select option to run, for example: install | uninstall | restart | run")
		os.Exit(0)
	}

	//Init all flags
	flags()

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
		log.Println("Bad runtime argument")
	}
}
