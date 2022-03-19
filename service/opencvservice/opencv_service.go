package opencvservice

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/rlaskowski/easymotion"
	"github.com/rlaskowski/easymotion/config"
)

type OpenCVService struct {
	captures     map[int]*Capture
	videosRecord map[int]*VideoRecord
}

func (OpenCVService) CreateService() *easymotion.ServiceInfo {
	return &easymotion.ServiceInfo{
		ID:        "service.opencv",
		Intstance: newCaptureService(),
	}
}

func newCaptureService() *OpenCVService {
	return &OpenCVService{
		captures:     make(map[int]*Capture),
		videosRecord: make(map[int]*VideoRecord),
	}
}

// Starting all process
//
// for example create active capture list
func (o *OpenCVService) Start() error {
	cap, err := OpenCapture(0)
	if err != nil {
		return err
	}

	o.captures[0] = cap

	return nil
}

// Stopping all active processes
func (o *OpenCVService) Stop() error {
	if err := o.StopRecording(0); err != nil {
		return err
	}

	if cap, ok := o.captures[0]; ok {
		return cap.Close()
	}

	return nil
}

// Finding capture by id
func (o *OpenCVService) Capture(id int) (*Capture, error) {
	cap, ok := o.captures[id]
	if !ok {
		return nil, fmt.Errorf("could not find capture %v", id)
	}

	return cap, nil
}

// Finding Video Record by capture id
func (o *OpenCVService) VideoRecord(id int) (*VideoRecord, error) {
	vr, ok := o.videosRecord[id]
	if !ok {
		return nil, fmt.Errorf("could not find video record, capture %v", id)
	}

	return vr, nil
}

// Stream video file
func (o *OpenCVService) Stream(capture *Capture) <-chan []byte {
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

// Starting recording by capture id
func (o *OpenCVService) StartRecording(id int) error {
	cap, err := o.Capture(id)
	if err != nil {
		return err
	}

	if _, err := o.VideoRecord(id); err == nil {
		return fmt.Errorf("video record is already exist, capture %v", id)
	}

	name := time.Now().Format("20060102_150405")
	videoPath := fmt.Sprintf("cam%d_%s.avi", id, name)

	vf, err := cap.VideoRecord(videoPath, "h264")
	if err != nil {
		return err
	}

	o.videosRecord[id] = vf

	go func() {
		for {
			if vf.Size() >= config.ToBytes(10) {
				if err := o.StopRecording(id); err != nil {
					log.Println(err)
					break
				}

				if err := o.StartRecording(id); err != nil {
					log.Println(err)
				}
				break
			}

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

// Stopping recording by capture id
func (o *OpenCVService) StopRecording(id int) error {
	vr, err := o.VideoRecord(id)
	if err != nil {
		return err
	}

	if err := vr.Close(); err != nil {
		return err
	}

	delete(o.videosRecord, id)

	return nil
}
