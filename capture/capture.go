package capture

import (
	"errors"
	"sync"

	"gocv.io/x/gocv"
)

type Capture struct {
	number   int
	campture *gocv.VideoCapture
	mat      gocv.Mat
	matPool  sync.Pool
}

func Open(number int) (*Capture, error) {
	capture, err := gocv.OpenVideoCapture(number)
	if err != nil {
		return nil, err
	}

	camera := &Capture{
		number:   number,
		campture: capture,
		mat:      gocv.NewMat(),
	}

	camera.matPool.New = func() interface{} {
		return gocv.NewMat()
	}

	return camera, nil
}

func (c *Capture) IsOpened() bool {
	return c.campture.IsOpened()
}

func (c *Capture) Close() error {
	if c.IsOpened() {
		return c.campture.Close()
	}
	return nil
}

func (c *Capture) Num() int {
	return c.number
}

func (c *Capture) VideoFile(name, codec string) (*VideoFile, error) {
	mat, err := c.readMat()
	if err != nil {
		return nil, err
	}

	writer, err := gocv.VideoWriterFile(name, codec, 30, mat.Cols(), mat.Rows(), true)
	if err != nil {
		return nil, err
	}

	v := &VideoFile{
		videoWriter: writer,
		capture:     c,
	}

	return v, nil

}

func (c *Capture) readMat() (gocv.Mat, error) {
	mat := c.mat

	if ok := c.campture.Read(&mat); !ok {
		return gocv.Mat{}, errors.New("nothing to read")
	}

	if c.mat.Empty() {
		return gocv.Mat{}, errors.New("empty mat")
	}

	return mat, nil
}

func (c *Capture) Read(b []byte) (n int, err error) {
	mat, err := c.readMat()
	if err != nil {
		return 0, err
	}

	/* if ok := c.campture.Read(&c.mat); !ok {
		return 0, nil
	}

	buff, err := gocv.IMEncode(".jpg", c.mat)
	if err != nil {
		return 0, err
	}

	if c.mat.Empty() {
		return 0, nil
	} */

	buff, err := gocv.IMEncode(".jpg", mat)
	if err != nil {
		return 0, err
	}

	n = copy(b, buff.GetBytes())
	buff.Close()

	return n, nil
}

func (c *Capture) acquireMat() *gocv.Mat {
	m := c.matPool.Get().(*gocv.Mat)
	return m
}
