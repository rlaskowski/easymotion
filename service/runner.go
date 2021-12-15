package service

import "github.com/rlaskowski/easymotion/db"

type Runner interface {
	Sqlite() *db.Sqlite
	CaptureService() *CaptureService
}

type ServiceRunner interface {
	Start() error
	Stop() error
}
