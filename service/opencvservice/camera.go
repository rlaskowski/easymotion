package opencvservice

import (
	"errors"
	"image"
	"image/color"
	"sync"
	"time"

	"github.com/rlaskowski/easymotion/config"
	"gocv.io/x/gocv"
)

type Camera struct {
	id      int
	capture *gocv.VideoCapture
	mat     gocv.Mat
	rwmu    sync.RWMutex
}

func OpenCamera(id int) (*Camera, error) {
	capture, err := gocv.OpenVideoCapture(id)
	if err != nil {
		return nil, err
	}

	camera := &Camera{
		id:      id,
		capture: capture,
		mat:     gocv.NewMat(),
	}

	return camera, nil
}

func (c *Camera) IsOpened() bool {
	return c.capture.IsOpened()
}

func (c *Camera) Close() error {
	if c.IsOpened() {
		return c.capture.Close()
	}
	return nil
}

// Returns camera number value registered in the system
func (c *Camera) ID() int {
	return c.id
}

/* func (c *Camera) VideoRecord(name, codec string) (*VideoRecord, error) {
	mat, err := c.readMat()

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
} */

// Reading current Mat value
func (c *Camera) Read(b []byte) (n int, err error) {
	c.rwmu.Lock()
	defer c.rwmu.Unlock()

	if ok := c.capture.Read(&c.mat); !ok {
		return 0, errors.New("unexpected error to read mat")
	}

	if c.mat.Empty() {
		return 0, nil
	}

	c.showDatetime()

	buff, err := gocv.IMEncode(".jpg", c.mat)
	if err != nil {
		return 0, err
	}

	n = copy(b, buff.GetBytes())
	buff.Close()

	return n, nil
}

// Starting recording to file system
/* func (c *Camera) StartRecord() error {
	recmux.RLock()

	if c.vrec != nil {
		recmux.Unlock()
		return fmt.Errorf("camera %d is still recording", c.id)
	}

	videoPath := filepath.Join(cmd.VideoPath, fmt.Sprintf("cam%d", c.id))

	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		if err := os.MkdirAll(videoPath, 0777); err != nil {
			recmux.Unlock()
			return fmt.Errorf("path: %s to store video file not exists", videoPath)
		}
	}

	name := fmt.Sprintf("%s.avi", time.Now().Format("20060102_150405"))
	videoPath = filepath.Join(videoPath, name)

	vr, err := c.VideoRecord(videoPath, "h264")
	if err != nil {
		recmux.Unlock()
		return err
	}

	c.vrec = vr

	recmux.Unlock()

	go func() {
		for vr.IsOpened() {
			if vr.Size() >= config.ToBytes(10) {
				if err := c.StopRecord(); err != nil {
					log.Println(err.Error())
					return
				}

				if err := c.StartRecord(); err != nil {
					log.Println(err.Error())
					return
				}
				break
			}

			c.rwmu.RLock()

			if c.mat.Empty() {
				c.rwmu.RUnlock()
				continue
			}

			err = vr.Write(c.mat)

			c.rwmu.RUnlock()

			if err != nil {
				if err == io.EOF {
					break
				}
				log.Println(err.Error())
				return
			}
		}

	}()

	return nil
}

// Stopping recording to file system
func (c *Camera) StopRecord() error {
	recmux.Lock()
	defer recmux.Unlock()

	if c.vrec == nil {
		return fmt.Errorf("camera %d nothing recording yet", c.id)
	}

	if err := c.vrec.Close(); err != nil {
		return err
	}

	c.vrec = nil

	return nil
} */

func (c *Camera) showDatetime() bool {
	if !config.OptionValue.CameraOption.Timeline {
		return false
	}

	white := color.RGBA{255, 255, 255, 0}
	t := time.Now().Format("02-01-2006 15:04:05")

	size := gocv.GetTextSize(t, gocv.FontHersheyPlain, 1, 1)
	pt := image.Pt((c.mat.Cols()-20)-(size.X), (c.mat.Rows()-20)-size.Y)

	gocv.PutText(&c.mat, t, pt, gocv.FontHersheyPlain, 1, white, 1)

	return true
}
