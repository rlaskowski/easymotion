package opencvservice

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"gocv.io/x/gocv"
)

const (
	JPEGCompress CompressType = "jpeg"
)

type CompressType string

type MatOption struct {
	Rows, Cols int
	Kind       int
	Data       []byte
}

func NewMatOption() *MatOption {
	return &MatOption{}
}

// Transferring gocv.Mat type to byte array
func (m *MatOption) ToBytes(cam *Camera) ([]byte, error) {
	mat, err := cam.ReadMat()
	if err != nil {
		return nil, err
	}

	b := mat.ToBytes()

	m.Rows = mat.Rows()
	m.Cols = mat.Cols()
	m.Kind = int(mat.Type())
	m.Data = b

	w := &bytes.Buffer{}
	enc := gob.NewEncoder(w)

	if err := enc.Encode(m); err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

// Creating gocv.Mat type from byte array
func (m *MatOption) Mat(data []byte) (gocv.Mat, error) {
	r := bytes.NewReader(data)
	dec := gob.NewDecoder(r)

	if err := dec.Decode(&m); err != nil {
		return gocv.Mat{}, err
	}

	mfb, err := gocv.NewMatFromBytes(int(m.Rows), int(m.Cols), gocv.MatType(m.Kind), m.Data)
	if err != nil {
		return gocv.Mat{}, err
	}

	return mfb, nil
}

func (m *MatOption) MatCompress(data []byte, c CompressType) ([]byte, error) {
	var fileExt gocv.FileExt

	switch c {
	case JPEGCompress:
		fileExt = gocv.JPEGFileExt
	default:
		return nil, fmt.Errorf("compress type '%s' not found", c)
	}

	mat, err := m.Mat(data)
	if err != nil {
		return nil, err
	}

	b, err := gocv.IMEncode(fileExt, mat)
	if err != nil {
		return nil, err
	}

	defer b.Close()

	return b.GetBytes(), nil
}
