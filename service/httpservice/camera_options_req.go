package httpservice

type CameraOptionsReq struct {
	// Custom name
	Name string `json:"name"`

	// Unique number in system
	CameraID int `json:"camera_id"`

	// Enable or disable auto recording after system start
	Autorec bool `json:"auto_recording"`

	// Showing date and time on all video
	Timeline bool `json:"timeline"`
}
