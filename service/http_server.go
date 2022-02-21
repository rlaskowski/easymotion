package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type HttpServer struct {
	cancel  context.CancelFunc
	context context.Context
	echo    *echo.Echo
	mutex   *sync.Mutex
}

func (HttpServer) CreateService() *ServiceInfo {
	return &ServiceInfo{
		ID:        "service.http.server",
		Intstance: newHttpServer(),
	}
}

func newHttpServer() *HttpServer {
	ctx, cancel := context.WithCancel(context.Background())

	return &HttpServer{
		cancel:  cancel,
		context: ctx,
		echo:    echo.New(),
		mutex:   &sync.Mutex{},
	}
}

func (h *HttpServer) prepareEndpoints() {
	h.echo.GET("/stream/:captureID", h.Stream)
	h.echo.POST("/capture/:captureID/recording/start", h.StartRecording)
	h.echo.POST("/capture/:captureID/recording/stop", h.StopRecording)
	h.echo.POST("/user/create", h.CreateUser)
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

func (h *HttpServer) StartRecording(c echo.Context) error {
	captureID, err := strconv.Atoi(c.Param("captureID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"capture ID problem": err.Error(),
		})
	}

	service, err := h.captureService()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"capture error": err.Error(),
		})
	}

	err = service.StartRecording(captureID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"video record problem": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, "Recording was started")
}

func (h *HttpServer) StopRecording(c echo.Context) error {
	captureID, err := strconv.Atoi(c.Param("captureID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"capture ID problem": err.Error(),
		})
	}

	service, err := h.captureService()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"capture error": err.Error(),
		})
	}

	err = service.StopRecording(captureID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"video record problem": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, "Recording was stopped")
}

func (h *HttpServer) Stream(c echo.Context) error {
	captureID, err := strconv.Atoi(c.Param("captureID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"capture ID problem": err.Error(),
		})
	}

	service, err := h.captureService()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"capture error": err.Error(),
		})
	}

	capture, err := service.Capture(captureID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"capture error": err.Error(),
		})
	}

	boundary := "STREAMCAMERA"

	c.Response().Header().Set("Content-Type", fmt.Sprintf("multipart/x-mixed-replace; boundary=%s", boundary))
	c.Response().WriteHeader(http.StatusOK)

	for {
		bch := service.Stream(capture)
		buff := <-bch

		b := bytes.NewBuffer(buff)

		mw := multipart.NewWriter(c.Response())
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

		c.Response().Flush()

	}

	return nil
}

func (h *HttpServer) CreateUser(c echo.Context) error {
	return nil
}

//Returns capture service after mapping from ServiceInfo struct
func (h *HttpServer) captureService() (*CaptureService, error) {
	service, err := GetService("service.capture")
	if err != nil {
		return nil, err
	}

	capture, ok := service.Intstance.(*CaptureService)
	if !ok {
		return nil, errors.New("bad mapping from ServiceInfo to CaptureService")
	}

	return capture, nil
}
