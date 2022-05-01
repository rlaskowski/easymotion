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
	h.echo.GET("/stream/:cameraID", h.Stream)
	h.echo.POST("/camera/:cameraID/recording/start", h.StartRecording)
	h.echo.POST("/camera/:cameraID/recording/stop", h.StopRecording)
	h.echo.POST("/user/create", h.CreateUser)
	h.echo.GET("/user", h.User)
	h.echo.POST("/camera/options/create", h.CreateOptions)
}

//Creating system user
func (h *HttpHandler) CreateUser(c echo.Context) error {
	parm := c.FormValue("user")

	user := &UserRequest{}
	if err := json.Unmarshal([]byte(parm), user); err != nil {
		return h.responseErr(c, UserResponseErr)
	}

	app := h.application()

	if err := app.CreateUser(user.Name, user.Email, user.Password); err != nil {
		return h.responseErr(c, NewUserErr)
	}

	return c.NoContent(http.StatusCreated)
}

func (h *HttpHandler) User(c echo.Context) error {
	app := h.application()

	users, err := app.Users()

	if err != nil {
		return h.responseErr(c, ResourcesErr)
	}

	return c.JSON(http.StatusOK, users)
}

func (h *HttpHandler) CreateOptions(c echo.Context) error {
	optReq := c.FormValue("option")

	options := &CameraOptionsReq{}

	if err := json.Unmarshal([]byte(optReq), options); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"camera option problem": err.Error(),
		})
	}

	app := h.application()

	if err := app.CreateOptions(options.CameraID, options.Name); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"camera option problem": err.Error(),
		})
	}

	return c.NoContent(http.StatusCreated)
}

func (h *HttpHandler) StartRecording(c echo.Context) error {
	cameraID, err := strconv.Atoi(c.Param("cameraID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"camera problem": err.Error(),
		})
	}

	app := h.application()

	if app.StartRecord(cameraID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"video record problem": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, "Recording was started")
}

func (h *HttpHandler) StopRecording(c echo.Context) error {
	cameraID, err := strconv.Atoi(c.Param("cameraID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"camera problem": err.Error(),
		})
	}

	app := h.application()

	if err := app.StopRecord(cameraID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"video record problem": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, "Recording was stopped")
}

func (h *HttpHandler) Stream(c echo.Context) error {
	cameraID, err := strconv.Atoi(c.Param("cameraID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"camera problem": err.Error(),
		})
	}

	app := h.application()

	boundary := "STREAMCAMERA"

	c.Response().Header().Set("Content-Type", fmt.Sprintf("multipart/x-mixed-replace; boundary=%s", boundary))
	c.Response().WriteHeader(http.StatusOK)

	for {
		r, _ := app.ReadBytes(cameraID)
		if err != nil {
			continue
		}

		buff := bytes.NewBuffer(r)

		mw := multipart.NewWriter(c.Response())
		header := make(textproto.MIMEHeader)

		mw.SetBoundary(boundary)

		header.Set("Content-Type", "image/jpeg")
		header.Set("Content-Length", fmt.Sprintf("%d", buff.Len()))

		w, err := mw.CreatePart(header)
		if err != nil {
			break
		}

		_, err = io.Copy(w, buff)
		if err != nil {
			break
		}

		c.Response().Flush()

	}

	return nil
}

func (h *HttpHandler) application() *app.App {
	app := h.appPool.Get().(*app.App)
	defer h.appPool.Put(app)

	return app
}

// Returns Error type according errors list definition
func (h *HttpHandler) responseErr(ctx echo.Context, err Error) error {
	return ctx.JSON(err.Code, err)
}
