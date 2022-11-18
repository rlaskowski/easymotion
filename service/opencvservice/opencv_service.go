package opencvservice

import (
	"log"

	"github.com/rlaskowski/manage"
)

type OpenCVService struct {
	camera *Camera
}

func (OpenCVService) CreateService() *manage.ServiceInfo {
	return &manage.ServiceInfo{
		ID:        "service.opencv",
		Priority:  1,
		Intstance: newCaptureService(),
	}
}

func newCaptureService() *OpenCVService {
	return &OpenCVService{}
}

// Starting all process
func (o *OpenCVService) Start() error {
	log.Println("starting openservice")
	cam, err := OpenCamera(0)
	if err != nil {
		return err
	}

	o.camera = cam

	return nil
}

// Stopping all active processes
func (o *OpenCVService) Stop() error {
	log.Println("stopping openservice")
	return o.camera.Close()
}

// Returns actual system camera instance
func (o *OpenCVService) Camera() *Camera {
	return o.camera
}
