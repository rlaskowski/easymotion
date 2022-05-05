package opencvservice

import (
	"errors"
	"sync"

	"github.com/rlaskowski/easymotion"
)

type OpenCVService struct {
	rwmu    *sync.RWMutex
	cameras map[int]*Camera
}

func (OpenCVService) CreateService() *easymotion.ServiceInfo {
	return &easymotion.ServiceInfo{
		ID:        "service.opencv",
		Intstance: newCaptureService(),
	}
}

func newCaptureService() *OpenCVService {
	return &OpenCVService{
		cameras: make(map[int]*Camera),
	}
}

// Starting all process
func (o *OpenCVService) Start() error {
	/* cam, err := OpenCamera(0)
	if err != nil {
		return err
	}

	o.cameras[0] = cam */

	return nil
}

// Stopping all active processes
func (o *OpenCVService) Stop() error {
	for _, camera := range o.cameras {
		o.rwmu.Lock()

		if err := camera.Close(); err != nil {
			return err
		}

		o.rwmu.Unlock()
	}

	return nil
}

// Finding camera instance in service list
func (o *OpenCVService) Camera(id int) (*Camera, error) {
	o.rwmu.RLock()
	defer o.rwmu.RUnlock()

	camera, ok := o.cameras[id]

	if !ok {
		return nil, errors.New("camera instance not found")
	}

	return camera, nil
}

// Creating new camera instance
func (o *OpenCVService) CreateCamera(id int, options CameraOptions) (*Camera, error) {
	o.rwmu.Lock()
	defer o.rwmu.Unlock()

	cam, err := OpenCamera(id, options)

	if err != nil {
		return nil, err
	}

	o.cameras[id] = cam

	return cam, nil
}
