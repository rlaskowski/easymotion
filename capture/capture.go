package capture

import (
	"errors"

	"gocv.io/x/gocv"
)

type Capture struct {
	number   int
	campture *gocv.VideoCapture
	mat      gocv.Mat
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

func (c *Capture) WriteToFile(name string) error {
	return errors.New("not yet implemented")
}

func (c *Capture) Read(b []byte) (n int, err error) {
	if ok := c.campture.Read(&c.mat); !ok {
		return 0, nil
	}

	buff, err := gocv.IMEncode(".jpg", c.mat)
	if err != nil {
		return 0, err
	}

	if c.mat.Empty() {
		return 0, nil
	}

	n = copy(b, buff.GetBytes())
	buff.Close()

	return n, nil
}
