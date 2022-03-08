package httpservice

import (
	"context"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rlaskowski/easymotion"
)

type HttpServer struct {
	cancel  context.CancelFunc
	context context.Context
	echo    *echo.Echo
	handler *HttpHandler
}

func (HttpServer) CreateService() *easymotion.ServiceInfo {
	return &easymotion.ServiceInfo{
		ID:        "service.http.server",
		Intstance: newHttpServer(),
	}
}

func newHttpServer() *HttpServer {
	ctx, cancel := context.WithCancel(context.Background())

	echo := echo.New()

	return &HttpServer{
		cancel:  cancel,
		context: ctx,
		echo:    echo,
		handler: NewHttpHandler(echo),
	}
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
	h.handler.CreateEndpoints()

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
