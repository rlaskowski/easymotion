package service

import (
	"errors"

	"github.com/rlaskowski/easymotion/camera"
)

type CameraService struct {
	cameras map[int]*camera.Camera
}

func NewCameraService() *CameraService {
	return &CameraService{
		cameras: make(map[int]*camera.Camera),
	}
}

func (c *CameraService) Start() error {
	return errors.New("not yet implemented")
}

func (c *CameraService) Stop() error {
	return errors.New("not yet implemented")
}
