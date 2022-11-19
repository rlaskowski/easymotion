package easymotion

import (
	"errors"
)

type DeviceCommand struct {
	// Showing information data on the video
	VideoInfo bool `json:"cam_info"`
}

func (d DeviceCommand) Do() error {
	return errors.New("method not yet implemented")
}
