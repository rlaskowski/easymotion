package opencvservice

import (
	"fmt"
	"io"
	"os"
	"path"
	"sync"
	"time"

	"github.com/rlaskowski/easymotion/config"
	"gocv.io/x/gocv"
)

var (
	recmux sync.RWMutex
)

type VideoRecord struct {
	name        string
	videoWriter *gocv.VideoWriter
}

func NewVideoRecord() *VideoRecord {
	name := fmt.Sprintf("cam%d_%s.avi", 0, time.Now().Format("20060102_150405"))
	path := path.Join(config.WorkingDirectory(), "data", name)

	return &VideoRecord{
		name: path,
	}
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
