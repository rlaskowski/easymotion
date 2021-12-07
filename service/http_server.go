package service

import (
	"context"
	"log"

	"github.com/labstack/echo/v4"
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
	h.echo.HideBanner = true
}

func (h *HttpServer) Start() error {
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

	if err := h.echo.Close(); err != nil {
		return err
	}

	return nil
}
