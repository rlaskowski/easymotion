package service

import "github.com/rlaskowski/easymotion/db"

type Runner interface {
	SqliteDB() *db.SqliteDB
	CameraService() *CameraService
}

type ServiceRunner interface {
	Start() error
	Stop() error
}
