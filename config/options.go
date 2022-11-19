package config

var (
	OptionValue = DefaultOptions
)

type Options struct {
	// RabbitMQ service address
	MQAddress string `json:"mq_address"`

	// Options for each camera
	CameraOption CameraOptions `json:"camera_options"`

	// ID of this service
	ServiceID string `json:"-"`

	// Port where http server is listening
	HTTPServerPort int `json:"httpserver_port"`
}

type CameraOptions struct {
	// Enable or disable auto recording after system start
	Autorec bool `json:"auto_recording"`

	// Showing date and time on all video
	Timeline bool `json:"timeline"`
}

var DefaultOptions = Options{
	MQAddress: hubURL(),
	CameraOption: CameraOptions{
		Autorec:  false,
		Timeline: true,
	},
	ServiceID:      "service.easymotion",
	HTTPServerPort: 9090,
}
