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

func (c *CaptureService) Stream(capture *capture.Capture) <-chan []byte {
	imgch := make(chan []byte, 100)

	go func() {
		buff := make([]byte, 1024*1024)

		_, err := capture.Read(buff)
		if err != nil {
			imgch <- nil
		}

		imgch <- buff
	}()

	return imgch
}

func (c *CaptureService) WriteFile(cap *capture.Capture) error {
	name := "example"
	videoPath := fmt.Sprintf("%s.avi", name)

	vf, err := cap.VideoFile(videoPath, "h264")
	if err != nil {
		return err
	}

	go func() {
		for {
			err := vf.Write()
			if err != nil {
				break
			}
		}
	}()

	return nil
}
