package easymotion

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/rlaskowski/easymotion/service/opencvservice"
	"github.com/rlaskowski/easymotion/service/queueservice"
	"github.com/rlaskowski/manage"
	"github.com/rlaskowski/manage/rabbitmq"
)

type Device struct {
	opencv    *opencvservice.OpenCVService
	mqservice *queueservice.RabbitMQService
}

func NewRunner() Runner {
	return newDevice()
}

func newDevice() *Device {
	return &Device{}
}

func (d *Device) Run() error {
	if err := d.services(); err != nil {
		return err
	}

	buff := &bytes.Buffer{}

	for {
		c := d.opencv.Camera()

		n, err := io.Copy(buff, c)
		if err != nil {
			log.Printf("reading camera error: %s", err.Error())
			break
		}

		if n > 0 {
			msg := rabbitmq.Message{
				ContentType: "image/jpg",
				Body:        buff.Bytes(),
			}

			if err := d.mqservice.Publish(context.Background(), msg); err != nil {
				log.Printf("rabbitmq publishing error: %s", err.Error())
			}
		}

	}

	return nil
}

func (d *Device) Close() error {
	return nil
}

// Initializing services
func (d *Device) services() error {
	mqservice, err := manage.GetService("service.rabbitmq")
	if err != nil {
		return fmt.Errorf("service.rabbitmq instance error: %s", err.Error())
	}

	mqinstance := mqservice.Intstance

	mq, ok := mqinstance.(*queueservice.RabbitMQService)
	if !ok {
		return errors.New("rabbitmq service assertion problem")
	}

	opencv, err := manage.GetService("service.opencv")
	if err != nil {
		return fmt.Errorf("service.opencv instance error: %s", err.Error())
	}

	ocvinstance := opencv.Intstance

	ocv, ok := ocvinstance.(*opencvservice.OpenCVService)
	if !ok {
		return errors.New("opencv service assertion problem")
	}

	d.mqservice = mq
	d.opencv = ocv

	return nil
}
