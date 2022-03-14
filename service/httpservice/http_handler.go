package httpservice

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rlaskowski/easymotion"
	"github.com/rlaskowski/easymotion/service/opencvservice"
)

type HttpHandler struct {
	echo *echo.Echo
}

func NewHttpHandler(echo *echo.Echo) *HttpHandler {
	return &HttpHandler{echo}
}

//Creating endpoints list
func (h *HttpHandler) CreateEndpoints() {
	h.echo.GET("/stream/:captureID", h.Stream)
	h.echo.POST("/capture/:captureID/recording/start", h.StartRecording)
	h.echo.POST("/capture/:captureID/recording/stop", h.StopRecording)
	h.echo.POST("/user/create", h.CreateUser)
}

//Creating system user
func (h *HttpHandler) CreateUser(c echo.Context) error {
	c.NoContent(http.StatusCreated)
	return nil
}

func (h *HttpHandler) StartRecording(c echo.Context) error {
	captureID, err := strconv.Atoi(c.Param("captureID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"capture ID problem": err.Error(),
		})
	}

	service, err := h.opencvService()
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

func (h *HttpHandler) StopRecording(c echo.Context) error {
	captureID, err := strconv.Atoi(c.Param("captureID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"capture ID problem": err.Error(),
		})
	}

	service, err := h.opencvService()
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

func (h *HttpHandler) Stream(c echo.Context) error {
	captureID, err := strconv.Atoi(c.Param("captureID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"capture ID problem": err.Error(),
		})
	}

	service, err := h.opencvService()
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

//Returns capture service after mapping from ServiceInfo struct
func (h *HttpHandler) opencvService() (*opencvservice.OpenCVService, error) {
	service, err := easymotion.GetService("service.opencv")
	if err != nil {
		return nil, err
	}

	opencvSrv, ok := service.Intstance.(*opencvservice.OpenCVService)
	if !ok {
		return nil, errors.New("bad mapping from ServiceInfo to OpenCVService")
	}

	return opencvSrv, nil
}
