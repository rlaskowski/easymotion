package dbservice

type CameraOptions struct {
	ID int64 `json:"-"`

	// Custom name
	Name string `json:"name"`

	// Unique number in system
	CameraID int `json:"-"`

	// Enable or disable auto recording after system start
	Autorec bool `json:"auto_recording"`

	// Showing date and time on all video
	Timeline bool `json:"timeline"`
}
