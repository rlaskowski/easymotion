package capture

import (
	"errors"
	"log"

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

func (c *Capture) Num() int {
	return c.number
}

func (c *Capture) VideoRecord(name, codec string) (*VideoRecord, error) {
	mat := <-c.readMat()

	if mat.Empty() {
		return nil, errors.New("empty mat")
	}

	writer, err := gocv.VideoWriterFile(name, codec, 30, mat.Cols(), mat.Rows(), true)
	if err != nil {
		return nil, err
	}

	v := &VideoRecord{
		videoWriter: writer,
		capture:     c,
	}

	return v, nil
}

func (c *Capture) readMat() <-chan gocv.Mat {
	match := make(chan gocv.Mat)

	go func() {

		if ok := c.campture.Read(&c.mat); !ok {
			log.Println("nothing to read from capture")
		}

		match <- c.mat

		close(match)

	}()

	return match
}

func (c *Capture) Read(b []byte) (n int, err error) {
	mat := <-c.readMat()

	if mat.Empty() {
		return 0, nil
	}

	buff, err := gocv.IMEncode(".jpg", mat)
	if err != nil {
		return 0, err
	}

	n = copy(b, buff.GetBytes())
	buff.Close()

	return n, nil
}
