package service

import "github.com/rlaskowski/easymotion/capture"

type CaptureService struct {
	cameras map[int]*capture.Capture
}

func NewCaptureService() *CaptureService {
	return &CaptureService{
		cameras: make(map[int]*capture.Capture),
	}
}

func (c *CaptureService) Start() error {
	cam, err := capture.Open(0)
	if err != nil {
		return err
	}

	c.cameras[0] = cam

	return nil
}

func (c *CaptureService) Stop() error {
	if cam, ok := c.cameras[0]; ok {
		return cam.Close()
	}

	return nil
}
