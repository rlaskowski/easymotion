package service

import (
	"fmt"
	"io"
	"log"

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

	if err := c.WriteFile(); err != nil {
		return err
	}

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

func (c *CaptureService) WriteFile() error {
	cap, err := c.Capture(0)
	if err != nil {
		return err
	}

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
				if err == io.EOF {
					break
				}
				log.Println(err)
			}
		}
	}()

	return nil
}
