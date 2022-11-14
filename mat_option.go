package easymotion

import (
	"unsafe"

	"github.com/rlaskowski/easymotion/service/opencvservice"
)

type MatOption struct {
	rows, cols uint8
	kind       uint8
	data       []byte
}

func MatToBytes(cam *opencvservice.Camera) ([]byte, error) {
	mat, err := cam.ReadMat()
	if err != nil {
		return nil, err
	}

	b := mat.ToBytes()

	moption := &MatOption{
		rows: uint8(mat.Rows()),
		cols: uint8(mat.Cols()),
		kind: uint8(mat.Type()),
		data: b,
	}

	buff := *(*[]byte)(unsafe.Pointer(&moption))

	/* buff := &bytes.Buffer{}
	if err := binary.Write(buff, binary.BigEndian, moption); err != nil {
		return nil, err
	} */

	return buff, nil
}
