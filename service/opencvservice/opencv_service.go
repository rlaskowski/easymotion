package opencvservice

import (
	"sync"

	"github.com/rlaskowski/easymotion"
)

type OpenCVService struct {
	rmu     *sync.RWMutex
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
		rmu:     &sync.RWMutex{},
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
		o.rmu.RLock()

		if err := camera.Close(); err != nil {
			return err
		}

		o.rmu.RUnlock()
	}

	return nil
}

// Finding or creating camera instance
func (o *OpenCVService) GetOrCreate(id int, options CameraOptions) (*Camera, error) {
	o.rmu.RLock()
	defer o.rmu.RUnlock()

	if cam, ok := o.cameras[id]; ok {
		return cam, nil
	}

	cam, err := OpenCamera(id, options)

	if err != nil {
		return nil, err
	}

	o.cameras[id] = cam

	return cam, nil
}
