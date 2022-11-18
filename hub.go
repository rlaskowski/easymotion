//go:build !device

package easymotion

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/rlaskowski/easymotion/config"
	"github.com/rlaskowski/easymotion/service/opencvservice"
	"github.com/rlaskowski/easymotion/service/queueservice"
	"github.com/rlaskowski/manage"
)

type Hub struct {
	mqservice *queueservice.RabbitMQService
}

func NewRunner() Runner {
	return newHub()
}

func newHub() *Hub {
	return &Hub{}
}

func (h *Hub) RegisterServices() error {
	manage.RegisterService(&queueservice.RabbitMQService{})

	return nil
}

func (h *Hub) Run() error {
	log.Println("running hub")

	if err := h.services(); err != nil {
		return err
	}

	sub, err := h.mqservice.Subscribe(context.Background())
	if err != nil {

		return err
	}

	var rec *opencvservice.VideoRecord

	go func() {
		for msg := range sub {
			//mat := *(*gocv.Mat)(unsafe.Pointer(&msg.Body))
			mat, err := MatFromBytes(msg.Body)
			if err != nil {
				log.Fatal(err)
				break
			}

			if rec == nil || rec.Size() >= config.ToBytes(10) {
				rec, err = opencvservice.OpenVideoRecord(mat)
				if err != nil {
					log.Fatal(err)
					break
				}
				continue
			}

			if err := rec.Write(mat); err != nil {
				continue
			}
		}
		rec.Close()
	}()

	return nil
}

func (h *Hub) Close() error {
	return errors.New("not yet implemented")
}

func (h *Hub) services() error {
	mqservice, err := manage.GetService("service.rabbitmq")
	if err != nil {
		return fmt.Errorf("service.rabbitmq instance error: %s", err.Error())
	}

	mqinstance := mqservice.Intstance

	mq, ok := mqinstance.(*queueservice.RabbitMQService)
	if !ok {
		return errors.New("rabbitmq service assertion problem")
	}

	h.mqservice = mq

	return nil
}
