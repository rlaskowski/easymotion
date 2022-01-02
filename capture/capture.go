package capture

import (
	"errors"
	"image"
	"image/color"
	"log"
	"time"

	"gocv.io/x/gocv"
)

var (
	DefaultCaptureOptions = CaptureOptions{
		timeline: true,
	}
)

type Capture struct {
	number   int
	campture *gocv.VideoCapture
	mat      gocv.Mat
	options  CaptureOptions
}

type CaptureOptions struct {
	//Showing date and time on all video
	timeline bool
}

func Open(number int) (*Capture, error) {
	return open(number, DefaultCaptureOptions)
}

func OpenWithOptions(number int, options CaptureOptions) (*Capture, error) {
	return open(number, options)
}

func open(number int, options CaptureOptions) (*Capture, error) {
	capture, err := gocv.OpenVideoCapture(number)
	if err != nil {
		return nil, err
	}

	camera := &Capture{
		number:   number,
		campture: capture,
		mat:      gocv.NewMat(),
		options:  options,
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
		name:        name,
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

		if c.options.timeline {
			c.showDatetime()
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

func (c *Capture) showDatetime() {
	white := color.RGBA{255, 255, 255, 0}
	t := time.Now().Format("02-01-2006 15:04:05")

	size := gocv.GetTextSize(t, gocv.FontHersheyPlain, 1, 1)
	//pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
	//pt := image.Pt(80+(80/2)-(size.X/2), 20)
	pt := image.Pt(70+(80/2)-(size.X/2), c.mat.Rows()-20)
	gocv.PutText(&c.mat, t, pt, gocv.FontHersheyPlain, 1, white, 1)
}
