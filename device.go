//go:build device

package easymotion

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/rlaskowski/easymotion/service/opencvservice"
	"github.com/rlaskowski/easymotion/service/queueservice"
	"github.com/rlaskowski/easymotion/service/storage"
	"github.com/rlaskowski/manage"
	"github.com/rlaskowski/manage/rabbitmq"
)

type Device struct {
	opencv     *opencvservice.OpenCVService
	mqservice  *queueservice.RabbitMQService
	ctx        context.Context
	cancel     context.CancelFunc
	matoptPool sync.Pool
}

func NewRunner() Runner {
	return newDevice()
}

func newDevice() *Device {
	ctx, cancel := context.WithCancel(context.Background())

	d := &Device{
		ctx:    ctx,
		cancel: cancel,
	}

	d.matoptPool.New = func() interface{} {
		return opencvservice.NewMatOption()
	}

	return d
}

func (d *Device) RegisterServices() {
	manage.RegisterService(&opencvservice.OpenCVService{})
	manage.RegisterService(&queueservice.RabbitMQService{})
	manage.RegisterService(&storage.SqliteService{})
}

func (d *Device) Run() error {
	log.Println("running device")

	if err := d.services(); err != nil {
		return err
	}

	go d.sendData()
	go d.receiveCommand()

	return nil
}

func (d *Device) Close() error {
	d.cancel()
	return nil
}

// Receiving executing commands from queue service
func (d *Device) receiveCommand() {

}

// Sending video capture to rabbitmq
func (d *Device) sendData() {
	c := d.opencv.Camera()

	for {
		matopt := d.matoptPool.Get().(*opencvservice.MatOption)
		defer d.matoptPool.Put(matopt)

		b, err := matopt.ToBytes(c)
		if err != nil {
			log.Printf("reading mat to bytes error: %s", err.Error())
			break
		}

		if len(b) > 0 {
			msg := rabbitmq.Message{
				ContentType: "application/octet-stream",
				Body:        b,
			}

			if err := d.mqservice.Publish(context.Background(), msg); err != nil {
				log.Printf("rabbitmq publishing error: %s", err.Error())
			}
		}
		select {
		case <-d.ctx.Done():
			log.Printf("running device error: %s", d.ctx.Err().Error())
			return
		default:
		}
	}
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
