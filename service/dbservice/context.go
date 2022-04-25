package dbservice

type Context interface {
	// Creating new users
	CreateUser(user *User) error

	// Returns user by email address
	UserByEmail(email string) (User, error)

	// Returns all user
	Users() ([]User, error)

	// Creating option for system camera
	CreateCamOption(options *CameraOptions) error

	// Returns camera option by camera id
	CameraOption(camID int) (CameraOptions, error)

	// Returns all camera options
	CameraOptions() ([]CameraOptions, error)
}
