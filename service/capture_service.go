package service

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/rlaskowski/easymotion/capture"
)

type CaptureService struct {
	captures     map[int]*capture.Capture
	videosRecord map[int]*capture.VideoRecord
}

func NewCaptureService() *CaptureService {
	return &CaptureService{
		captures:     make(map[int]*capture.Capture),
		videosRecord: make(map[int]*capture.VideoRecord),
	}
}

//Starting all process
//
//for example create active capture list
func (c *CaptureService) Start() error {
	cap, err := capture.Open(0)
	if err != nil {
		return err
	}

	c.captures[0] = cap

	return nil
}

//Stopping all active processes
func (c *CaptureService) Stop() error {
	if err := c.StopRecording(0); err != nil {
		return err
	}

	if cap, ok := c.captures[0]; ok {
		return cap.Close()
	}

	return nil
}

//Finding capture by id
func (c *CaptureService) Capture(id int) (*capture.Capture, error) {
	cap, ok := c.captures[id]
	if !ok {
		return nil, fmt.Errorf("could not find capture %v", id)
	}

	return cap, nil
}

//Finding Video Record by capture id
func (c *CaptureService) VideoRecord(id int) (*capture.VideoRecord, error) {
	vr, ok := c.videosRecord[id]
	if !ok {
		return nil, fmt.Errorf("could not find video record, capture %v", id)
	}

	return vr, nil
}

//Stream video file
func (c *CaptureService) Stream(capture *capture.Capture) <-chan []byte {
	imgch := make(chan []byte, 10)

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

//Starting recording by capture id
func (c *CaptureService) StartRecording(id int) error {
	cap, err := c.Capture(id)
	if err != nil {
		return err
	}

	if _, err := c.VideoRecord(id); err == nil {
		return fmt.Errorf("video record is already exist, capture %v", id)
	}

	name := time.Now().Format("2006-01-02_15-04-05")
	videoPath := fmt.Sprintf("cam%d_%s.avi", id, name)

	vf, err := cap.VideoRecord(videoPath, "h264")
	if err != nil {
		return err
	}

	c.videosRecord[id] = vf

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

//Stopping recording by capture id
func (c *CaptureService) StopRecording(id int) error {
	vr, err := c.VideoRecord(id)
	if err != nil {
		return err
	}

	if err := vr.Close(); err != nil {
		return err
	}

	delete(c.videosRecord, id)

	return nil
}
