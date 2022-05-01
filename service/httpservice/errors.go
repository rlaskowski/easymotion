package httpservice

import (
	"net/http"
)

var (
	ApplicationErr  = newError(http.StatusBadRequest, "couldn't find application context")
	UserResponseErr = newError(http.StatusBadRequest, "bad definition of user response")
	NewUserErr      = newError(http.StatusBadRequest, "user couldn't be added")
	ResourcesErr    = newError(http.StatusNotFound, "resources not found")
)

type Error struct {
	Name string `json:"error_name"`
	Code int    `json:"code"`
}

func newError(code int, name string) Error {
	return Error{Name: name, Code: code}
}
