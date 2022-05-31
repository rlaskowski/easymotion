package opencvservice

import (
	"errors"
	"sync"

	"github.com/rlaskowski/easymotion"
)

type OpenCVService struct {
	cammu   *sync.RWMutex
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
		cammu:   &sync.RWMutex{},
		cameras: make(map[int]*Camera),
	}
}

// Starting all process
func (o *OpenCVService) Start() error {
	return nil
}

// Stopping all active processes
func (o *OpenCVService) Stop() error {
	for _, camera := range o.cameras {
		o.cammu.Lock()

		if err := camera.Close(); err != nil {
			o.cammu.Unlock()
			return err
		}

		o.cammu.Unlock()
	}

	return nil
}

// Finding camera instance in service list
func (o *OpenCVService) Camera(id int) (*Camera, error) {
	o.cammu.RLock()
	defer o.cammu.RUnlock()

	camera, ok := o.cameras[id]

	if !ok {
		return nil, errors.New("camera instance not found")
	}

	return camera, nil
}

// Creating new camera instance
func (o *OpenCVService) CreateCamera(id int, options CameraOptions) (*Camera, error) {
	o.cammu.Lock()
	defer o.cammu.Unlock()

	cam, err := OpenCamera(id, options)

	if err != nil {
		return nil, err
	}

	go cam.ReadMat()

	o.cameras[id] = cam

	return cam, nil
}
