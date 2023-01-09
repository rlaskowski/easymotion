package opencvservice

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/rlaskowski/easymotion/config"
	"gocv.io/x/gocv"
)

type VideoRecord struct {
	path        string
	videoWriter *gocv.VideoWriter
}

func OpenVideoRecord(mat gocv.Mat) (*VideoRecord, error) {
	path := config.RecordsPath()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0777); err != nil {
			return nil, fmt.Errorf("path: %s to store video file not exists", path)
		}
	}

	name := fmt.Sprintf("cam%d_%s.avi", 0, time.Now().Format("20060102_150405"))
	path = filepath.Join(path, name)

	writer, err := gocv.VideoWriterFile(path, "h264", 30, mat.Cols(), mat.Rows(), true)
	if err != nil {
		return nil, err
	}

	return &VideoRecord{
		path:        path,
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
	fi, err := os.Stat(v.path)
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
