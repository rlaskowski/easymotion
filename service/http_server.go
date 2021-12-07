package service

import (
	"context"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type HttpServer struct {
	cancel  context.CancelFunc
	context context.Context
	echo    *echo.Echo
}

func NewHttpServer() *HttpServer {
	ctx, cancel := context.WithCancel(context.Background())

	return &HttpServer{
		cancel:  cancel,
		context: ctx,
		echo:    echo.New(),
	}
}

func (h *HttpServer) prepareEndpoints() {

}

func (h *HttpServer) configure() {
	h.echo.HideBanner = true
	h.echo.HidePort = true
	h.echo.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `Method: ${method}, Path: ${path}, Remote IP: ${remote_ip}, Status: ${status}`,
	}))
}

func (h *HttpServer) Start() error {
	log.Println("Starting Http Server")

	h.configure()
	h.prepareEndpoints()

	go func() {
		if err := h.echo.Start(":9090"); err != nil {
			log.Fatalf("could not start http server due to: %s", err.Error())
		}
	}()

	return nil
}

func (h *HttpServer) Stop() error {
	h.cancel()

	log.Println("Stopping Http Server")

	if err := h.echo.Close(); err != nil {
		return err
	}

	return nil
}
