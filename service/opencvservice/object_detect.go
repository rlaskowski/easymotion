package opencvservice

import (
	"errors"

	"gocv.io/x/gocv"
)

type ObjectDetect struct {
	classifier gocv.CascadeClassifier
	capture    *Capture
}

func (o *ObjectDetect) Face() error {
	return errors.New("not yet implemented")
}

func (o *ObjectDetect) Close() error {
	return o.classifier.Close()
}
