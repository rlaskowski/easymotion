package httpservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/rlaskowski/easymotion/app"
)

type HttpHandler struct {
	echo    *echo.Echo
	appPool sync.Pool
}

func NewHttpHandler(echo *echo.Echo) *HttpHandler {
	return &HttpHandler{
		echo: echo,
		appPool: sync.Pool{
			New: func() interface{} {
				return app.NewApp()
			},
		},
	}
}

//Creating endpoints list
func (h *HttpHandler) CreateEndpoints() {
	h.echo.GET("/stream/:captureID", h.Stream)
	h.echo.POST("/capture/:captureID/recording/start", h.StartRecording)
	h.echo.POST("/capture/:captureID/recording/stop", h.StopRecording)
	h.echo.POST("/user/create", h.CreateUser)
	h.echo.GET("/user", h.User)
}

//Creating system user
func (h *HttpHandler) CreateUser(c echo.Context) error {
	parm := c.FormValue("user")

	user := &UserRequest{}
	if err := json.Unmarshal([]byte(parm), user); err != nil {
		return h.responseErr(c, UserResponseErr)
	}

	app := h.app()
	if err := app.CreateUser(user.Name, user.Email, user.Password); err != nil {
		return h.responseErr(c, NewUserErr)
	}
	c.NoContent(http.StatusCreated)

	return nil
}

func (h *HttpHandler) User(c echo.Context) error {
	app := h.app()
	users, err := app.Users()

	if err != nil {
		return h.responseErr(c, ResourcesErr)
	}

	return c.JSON(http.StatusOK, users)
}

func (h *HttpHandler) StartRecording(c echo.Context) error {
	captureID, err := strconv.Atoi(c.Param("captureID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"capture ID problem": err.Error(),
		})
	}

	app := h.app()
	if app == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"capture error": "bad application instance",
		})
	}

	err = app.OpenCVService().StartRecording(captureID)
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

	app := h.app()
	if app == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"capture error": "bad application instance",
		})
	}

	err = app.OpenCVService().StopRecording(captureID)
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

	app := h.app()
	if app == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"stream error": "bad application instance",
		})
	}

	capture, err := app.OpenCVService().Capture(captureID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"stream error": err.Error(),
		})
	}

	boundary := "STREAMCAMERA"

	c.Response().Header().Set("Content-Type", fmt.Sprintf("multipart/x-mixed-replace; boundary=%s", boundary))
	c.Response().WriteHeader(http.StatusOK)

	for {
		bch := app.OpenCVService().Stream(capture)
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

func (h *HttpHandler) app() *app.App {
	app := h.appPool.Get().(*app.App)
	defer h.appPool.Put(app)

	return app
}

// Returns Error type according errors list definition
func (h *HttpHandler) responseErr(ctx echo.Context, err Error) error {
	return ctx.JSON(err.Code, err)
}
