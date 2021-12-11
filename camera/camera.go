package camera

import (
	"errors"

	"gocv.io/x/gocv"
)

type Camera struct {
	number   int
	campture *gocv.VideoCapture
	mat      *gocv.Mat
}

func (c *Camera) Open(number int) (*Camera, error) {
	capture, err := gocv.OpenVideoCapture(number)
	if err != nil {
		return nil, err
	}

	camera := &Camera{
		number:   number,
		campture: capture,
	}

	return camera, nil
}

func (c *Camera) Close() error {
	return c.campture.Close()
}

func (c *Camera) WriteToFile(name string) error {
	return errors.New("not yet implemented")
}

func (c *Camera) Read(b []byte) error {
	mat := gocv.NewMat()
	defer mat.Close()

	for c.campture.IsOpened() {
		if ok := c.campture.Read(&mat); !ok {
			return nil
		}

		if mat.Empty() {
			continue
		}

		buff, err := gocv.IMEncode(".jpg", mat)
		if err != nil {
			return err
		}

		b = buff.GetBytes()
		buff.Close()
	}

	return nil
}
