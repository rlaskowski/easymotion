package service

import (
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
	cam, err := camera.Open(0)
	if err != nil {
		return err
	}

	c.cameras[0] = cam

	return nil
}

func (c *CameraService) Stop() error {
	if cam, ok := c.cameras[0]; ok {
		return cam.Close()
	}

	return nil
}
