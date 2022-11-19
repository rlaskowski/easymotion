package httpservice

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/rlaskowski/easymotion/config"
	"github.com/rlaskowski/manage"
)

type HttpServer struct {
	echo   *echo.Echo
	router Router
}

func (HttpServer) CreateService() *manage.ServiceInfo {
	return &manage.ServiceInfo{
		ID:        "service.http.server",
		Priority:  2,
		Intstance: newHttpServer(),
	}
}

func newHttpServer() *HttpServer {
	return &HttpServer{
		echo: echo.New(),
	}
}

func (h *HttpServer) configure() {
	h.echo.HideBanner = true
	h.echo.HidePort = true
	h.echo.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `Method: ${method}, Path: ${path}, Remote IP: ${remote_ip}, Status: ${status}`,
	}))

	h.echo.GET("/video/stream", h.Stream)

}

func (h *HttpServer) Start() error {
	log.Println("starting http server")

	h.configure()

	go func() {
		if err := h.echo.Start(fmt.Sprintf(":%v", config.DefaultOptions.HTTPServerPort)); err != nil {
			log.Fatalf("could not start http server due to: %s", err.Error())
		}
	}()

	return nil
}

func (h *HttpServer) Stop() error {
	log.Println("stopping http server")

	return h.echo.Close()
}

func (h *HttpServer) ResetRouter(router Router) {
	h.router = router
}

func (h *HttpServer) Stream(ctx echo.Context) error {
	boundary := "--stream--"

	ctx.Response().Header().Set("Content-Type", fmt.Sprintf("multipart/x-mixed-replace; boundary=%s", boundary))
	ctx.Response().WriteHeader(http.StatusOK)

	for {
		stream := h.router.StreamVideo()

		for s := range stream {
			if s.Err != nil {
				return ctx.JSON(http.StatusBadRequest, s.Err.Error())
			}

			b := bytes.NewBuffer(s.Data)

			mw := multipart.NewWriter(ctx.Response())
			header := make(textproto.MIMEHeader)

			mw.SetBoundary(boundary)

			header.Set("Content-Type", "image/jpeg")
			header.Set("Content-Length", fmt.Sprintf("%d", b.Len()))

			w, err := mw.CreatePart(header)
			if err != nil {
				break
			}

			_, err = io.Copy(w, b)
			if err != nil {
				break
			}

			ctx.Response().Flush()
		}

	}
}
