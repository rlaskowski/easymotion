package service

type Runner interface {
	CaptureService() *CaptureService
}

type ServiceRunner interface {
	Start() error
	Stop() error
}
