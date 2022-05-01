package opencvservice

import (
	"io"
	"os"

	"gocv.io/x/gocv"
)

var (
	actualRec = make(map[int]*VideoRecord)
)

type VideoRecord struct {
	name        string
	videoWriter *gocv.VideoWriter
}

func (v *VideoRecord) Write(mat gocv.Mat) error {
	if mat.Empty() {
		return nil
	}

	if v.IsOpened() {
		v.videoWriter.Write(mat)
	} else {
		return io.EOF
	}

	return nil
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
