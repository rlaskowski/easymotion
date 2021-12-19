package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rlaskowski/easymotion/service"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Please select option to run for example: install | uninstall | restart | run")
	}

	service, err := service.CreateSystemContext()
	if err != nil {
		fmt.Errorf("Unexpected error: %s", err.Error())
	}

	switch os.Args[1] {
	case "run":
		if err := service.RunService(); err != nil {
			log.Println(err)
		}
	case "install":
		if err := service.InstallService(); err != nil {
			log.Println(err)
		}
	case "uninstall":
		if err := service.UninstallService(); err != nil {
			log.Println(err)
		}
	case "restart":
		if err := service.RestartService(); err != nil {
			log.Println(err)
		}
	}
	/* service := service.NewService()
	service.Start() */
}
