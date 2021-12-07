package main

import (
	"github.com/rlaskowski/easymotion/service"
)

func main() {
	service := service.NewService()
	service.Start()
}
