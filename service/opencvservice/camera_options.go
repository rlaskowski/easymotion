package opencvservice

type CameraOptions struct {
	// Enable or disable auto recording after system start
	Autorec bool `json:"auto_recording"`

	// Showing date and time on all video
	Timeline bool `json:"timeline"`
}
