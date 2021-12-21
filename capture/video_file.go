package capture

import (
	"errors"
	"io"

	"gocv.io/x/gocv"
)

type VideoFile struct {
	videoWriter *gocv.VideoWriter
	capture     *Capture
}

func (v *VideoFile) Write() error {
	mat := <-v.capture.readMat()
	if mat.Empty() {
		return errors.New("empty mat")
	}

	if v.videoWriter.IsOpened() {
		v.videoWriter.Write(mat)
	} else {
		return io.EOF
	}

	return nil
}

func (v *VideoFile) Close() error {
	return v.videoWriter.Close()
}
