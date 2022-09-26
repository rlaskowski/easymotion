package config

var (
	OptionValue = DefaultOptions
)

type Options struct {
	// Options for each camera
	CameraOption CameraOptions `json:"camera_options"`

	// ID of this service
	ServiceID string `json:"-"`
}

type CameraOptions struct {
	// Enable or disable auto recording after system start
	Autorec bool `json:"auto_recording"`

	// Showing date and time on all video
	Timeline bool `json:"timeline"`
}

var DefaultOptions = Options{
	CameraOption: CameraOptions{
		Autorec:  false,
		Timeline: true,
	},
	ServiceID: "service.easymotion",
}
