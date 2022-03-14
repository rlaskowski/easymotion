package main

import (
	"log"

	"github.com/rlaskowski/easymotion"
	"github.com/rlaskowski/easymotion/cmd"
	"github.com/rlaskowski/easymotion/config"
)

func main() {
	systemService, err := easymotion.CreateSystemService(config.ProjectPath())
	if err != nil {
		log.Fatalf("Unexpected error: %s", err.Error())
	}

	cmd.RunCommand(systemService)
}
