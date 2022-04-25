package opencvservice

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rlaskowski/easymotion/cmd"
	"github.com/rlaskowski/easymotion/config"
	"gocv.io/x/gocv"
)

type Camera struct {
	id       int
	campture *gocv.VideoCapture
	mat      gocv.Mat
	mu       sync.Mutex
	options  CameraOptions
}

func OpenCamera(id int, options CameraOptions) (*Camera, error) {
	capture, err := gocv.OpenVideoCapture(id)
	if err != nil {
		return nil, err
	}

	camera := &Camera{
		id:       id,
		campture: capture,
		mat:      gocv.NewMat(),
		options:  options,
	}

	return camera, nil
}

func (c *Camera) IsOpened() bool {
	return c.campture.IsOpened()
}

func (c *Camera) Close() error {
	if c.IsOpened() {
		return c.campture.Close()
	}
	return nil
}

// Returns camera number value registered in the system
func (c *Camera) ID() int {
	return c.id
}

func (c *Camera) VideoRecord(name, codec string) (*VideoRecord, error) {
	if err := c.readMat(); err != nil {
		return nil, err
	}

	if c.mat.Empty() {
		return nil, errors.New("to write a video file empty mat is not acceptable")
	}

	writer, err := gocv.VideoWriterFile(name, codec, 30, c.mat.Cols(), c.mat.Rows(), true)
	if err != nil {
		return nil, err
	}

	v := &VideoRecord{
		name:        name,
		videoWriter: writer,
	}

	return v, nil
}

func (c *Camera) readMat() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ok := c.campture.Read(&c.mat); !ok {
		return errors.New("unexpected error to read mat")
	}

	c.showDatetime()

	return nil
}

// Reading current Mat value
func (c *Camera) Read(b []byte) (n int, err error) {
	err = c.readMat()

	if err != nil {
		return 0, err
	}

	if c.mat.Empty() {
		return 0, nil
	}

	buff, err := gocv.IMEncode(".jpg", c.mat)
	if err != nil {
		return 0, err
	}

	n = copy(b, buff.GetBytes())
	buff.Close()

	return n, nil
}

// Starting recording to file system
func (c *Camera) StartRecord() error {
	if _, err := os.Stat(cmd.VideoPath); os.IsNotExist(err) {
		return fmt.Errorf("path: %s to store video file not exists", cmd.VideoPath)
	}

	name := time.Now().Format("20060102_150405")
	videoPath := filepath.Join(cmd.VideoPath, fmt.Sprintf("cam%d", c.id), name)

	vr, err := c.VideoRecord(videoPath, "h264")
	if err != nil {
		return err
	}

	rec, ok := actualRec[c.id]

	if ok {
		return fmt.Errorf("camera %d is still recording", c.id)
	}

	actualRec[c.id] = rec

	for {
		if vr.Size() >= config.ToBytes(10) {
			if err := c.StopRecord(); err != nil {
				return err
			}

			if err := c.StartRecord(); err != nil {
				return err
			}
			break
		}

		if err := c.readMat(); err != nil {
			if err := c.StopRecord(); err != nil {
				return err
			}

			return err
		}

		err := vr.Write(c.mat)

		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}

	return nil
}

// Stopping recording to file system
func (c *Camera) StopRecord() error {
	rec, ok := actualRec[c.id]

	if !ok {
		return fmt.Errorf("camera %d nothing recording yet", c.id)
	}

	if err := rec.Close(); err != nil {
		return err
	}

	delete(actualRec, c.id)

	return nil
}

func (c *Camera) showDatetime() bool {
	if !c.options.Timeline {
		return false
	}

	white := color.RGBA{255, 255, 255, 0}
	t := time.Now().Format("02-01-2006 15:04:05")

	size := gocv.GetTextSize(t, gocv.FontHersheyPlain, 1, 1)
	pt := image.Pt((c.mat.Cols()-20)-(size.X), (c.mat.Rows()-20)-size.Y)

	gocv.PutText(&c.mat, t, pt, gocv.FontHersheyPlain, 1, white, 1)

	return true
}
