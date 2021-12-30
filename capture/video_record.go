package capture

import (
	"errors"
	"io"
	"os"

	"gocv.io/x/gocv"
)

type VideoRecord struct {
	name        string
	videoWriter *gocv.VideoWriter
	capture     *Capture
}

func (v *VideoRecord) Write() error {
	mat := <-v.capture.readMat()
	if mat.Empty() {
		return errors.New("empty mat")
	}

	if v.IsOpened() {
		v.videoWriter.Write(mat)
	} else {
		return io.EOF
	}

	return nil
}

//Writes until data not across number of bytes declared in n
func (v *VideoRecord) WriteTo(n int64) error {
	return errors.New("not yet implemented")
}

func (v *VideoRecord) Size() int64 {
	fi, err := os.Stat(v.name)
	if err != nil {
		return 0
	}

	return fi.Size()
}

func (v *VideoRecord) IsOpened() bool {
	return v.videoWriter.IsOpened()
}

func (v *VideoRecord) Close() error {
	return v.videoWriter.Close()
}
