package opencvservice

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/rlaskowski/easymotion/config"
	"gocv.io/x/gocv"
)

type VideoRecord struct {
	name        string
	videoWriter *gocv.VideoWriter
}

func OpenVideoRecord(mat gocv.Mat) (*VideoRecord, error) {
	name := fmt.Sprintf("cam%d_%s.avi", 0, time.Now().Format("20060102_150405"))
	path := path.Join(config.WorkingDirectory(), name)

	writer, err := gocv.VideoWriterFile(name, "h264", 30, mat.Cols(), mat.Rows(), true)
	if err != nil {
		return nil, err
	}

	return &VideoRecord{
		name:        path,
		videoWriter: writer,
	}, nil
}

func (v *VideoRecord) Write(mat gocv.Mat) error {
	if mat.Empty() {
		return nil
	}

	if v.IsOpened() {
		if err := v.videoWriter.Write(mat); err != nil {
			return err
		}
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
