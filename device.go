//go:build device

package easymotion

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/rlaskowski/easymotion/config"
	"github.com/rlaskowski/easymotion/service/opencvservice"
	"github.com/rlaskowski/manage"
)

type Device struct {
	opencv *opencvservice.OpenCVService
	ctx    context.Context
	cancel context.CancelFunc
}

func NewRunner() Runner {
	return newDevice()
}

func newDevice() *Device {
	ctx, cancel := context.WithCancel(context.Background())

	return &Device{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (d *Device) RegisterServices() error {
	manage.RegisterService(&opencvservice.OpenCVService{})

	return nil
}

func (d *Device) Run() error {
	log.Println("running device")

	if err := d.services(); err != nil {
		return err
	}

	var rec *opencvservice.VideoRecord
	c := d.opencv.Camera()

	go func() {
		defer rec.Close()

		for {
			mat, err := c.ReadMat()
			if err != nil {
				log.Printf("read mat error: %s", err.Error())
				continue
			}

			if rec == nil || rec.Size() >= config.ToBytes(10) {
				rec, err = opencvservice.OpenVideoRecord(*mat)
				if err != nil {
					log.Fatal(err)
					break
				}
				continue
			}

			if err := rec.Write(*mat); err != nil {
				continue
			}

			select {
			case <-d.ctx.Done():
				return
			default:
			}
		}
	}()

	return nil
}

func (d *Device) Close() error {
	d.cancel()
	return nil
}

// Initializing services
func (d *Device) services() error {
	opencv, err := manage.GetService("service.opencv")
	if err != nil {
		return fmt.Errorf("service.opencv instance error: %s", err.Error())
	}

	ocvinstance := opencv.Intstance

	ocv, ok := ocvinstance.(*opencvservice.OpenCVService)
	if !ok {
		return errors.New("opencv service assertion problem")
	}

	d.opencv = ocv

	return nil
}
