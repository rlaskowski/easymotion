package service

import (
	"bytes"
	"context"
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
	Runner  Runner
	mutex   *sync.Mutex
}

func NewHttpServer(runner Runner) *HttpServer {
	ctx, cancel := context.WithCancel(context.Background())

	return &HttpServer{
		cancel:  cancel,
		context: ctx,
		echo:    echo.New(),
		Runner:  runner,
		mutex:   new(sync.Mutex),
	}
}

func (h *HttpServer) prepareEndpoints() {
	h.echo.GET("/stream/:captureID", h.Stream)
	h.echo.POST("/capture/record/:captureID", h.StartRecord)
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

func (h *HttpServer) StartRecord(c echo.Context) error {
	captureID, err := strconv.Atoi(c.Param("captureID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"capture ID problem": err.Error(),
		})
	}

	capture, err := h.Runner.CaptureService().Capture(captureID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"capture error": err.Error(),
		})
	}

	err = h.Runner.CaptureService().WriteFile(capture)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"video record problem": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, "Started recordin")
}

func (h *HttpServer) Stream(c echo.Context) error {
	captureID, err := strconv.Atoi(c.Param("captureID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"capture ID problem": err.Error(),
		})
	}

	capture, err := h.Runner.CaptureService().Capture(captureID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"capture error": err.Error(),
		})
	}

	//buff := make([]byte, 1024*1024)
	boundary := "STREAMCAMERA"

	c.Response().Header().Set("Content-Type", fmt.Sprintf("multipart/x-mixed-replace; boundary=%s", boundary))
	c.Response().WriteHeader(http.StatusOK)

	for {
		/* 	n, err := cam.Read(buff)
		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]string{
				"error": fmt.Sprintf("camera read problem: %s", err.Error()),
			})
			cam.Close()
			break
		} */

		bch := h.Runner.CaptureService().Stream(capture)
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
