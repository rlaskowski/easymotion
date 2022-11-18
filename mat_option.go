package easymotion

import (
	"bytes"
	"encoding/gob"

	"github.com/rlaskowski/easymotion/service/opencvservice"
	"gocv.io/x/gocv"
)

type MatOption struct {
	Rows, Cols int
	Kind       int
	Data       []byte
}

func MatToBytes(cam *opencvservice.Camera) ([]byte, error) {
	mat, err := cam.ReadMat()
	if err != nil {
		return nil, err
	}

	b := mat.ToBytes()

	moption := &MatOption{
		Rows: mat.Rows(),
		Cols: mat.Cols(),
		Kind: int(mat.Type()),
		Data: b,
	}
	w := &bytes.Buffer{}
	enc := gob.NewEncoder(w)
	if err := enc.Encode(moption); err != nil {
		return nil, err
	}

	//buff := *(*[]byte)(unsafe.Pointer(&moption))

	/* buff := &bytes.Buffer{}
	if err := binary.Write(buff, binary.BigEndian, moption); err != nil {
		return nil, err
	} */

	return w.Bytes(), nil
}

func MatFromBytes(data []byte) (gocv.Mat, error) {
	moption := &MatOption{}

	r := bytes.NewReader(data)
	dec := gob.NewDecoder(r)
	if err := dec.Decode(moption); err != nil {
		return gocv.Mat{}, err
	}

	mat, err := gocv.NewMatFromBytes(int(moption.Rows), int(moption.Cols), gocv.MatType(moption.Kind), moption.Data)
	if err != nil {
		return gocv.Mat{}, err
	}

	return mat, nil

}
