package capture

import (
	"gocv.io/x/gocv"
)

type VideoFile struct {
	videoWriter *gocv.VideoWriter
	capture     *Capture
}

func (v *VideoFile) Write() error {
	mat, err := v.capture.readMat()
	if err != nil {
		return err
	}

	if v.videoWriter.IsOpened() {
		v.videoWriter.Write(mat)
	}

	return nil
}

func (v *VideoFile) Close() error {
	return v.videoWriter.Close()
}
