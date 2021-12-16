package service

import (
	"fmt"

	"github.com/rlaskowski/easymotion/capture"
)

type CaptureService struct {
	captures map[int]*capture.Capture
}

func NewCaptureService() *CaptureService {
	return &CaptureService{
		captures: make(map[int]*capture.Capture),
	}
}

func (c *CaptureService) Start() error {
	cam, err := capture.Open(0)
	if err != nil {
		return err
	}

	c.captures[0] = cam

	return nil
}

func (c *CaptureService) Stop() error {
	if cam, ok := c.captures[0]; ok {
		return cam.Close()
	}

	return nil
}

//Finding capture by id
func (c *CaptureService) Capture(id int) (*capture.Capture, error) {
	capture, ok := c.captures[id]

	if !ok {
		return nil, fmt.Errorf("capture ID %d not found", id)
	}

	return capture, nil
}
