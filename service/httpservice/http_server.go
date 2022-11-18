package httpservice

import (
	"log"

	"github.com/rlaskowski/manage"
)

type HttpServer struct {
}

func (HttpServer) CreateService() *manage.ServiceInfo {
	return &manage.ServiceInfo{
		ID:        "service.http.server",
		Priority:  2,
		Intstance: newHttpServer(),
	}
}

func newHttpServer() *HttpServer {
	return &HttpServer{}
}

func (h *HttpServer) Start() error {
	log.Println("starting http server")

	return nil
}

func (h *HttpServer) Stop() error {
	log.Println("stopping http server")
	return nil
}
