//go:build !device

package easymotion

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/rlaskowski/easymotion/config"
	"github.com/rlaskowski/easymotion/service/httpservice"
	"github.com/rlaskowski/easymotion/service/opencvservice"
	"github.com/rlaskowski/easymotion/service/queueservice"
	"github.com/rlaskowski/manage"
)

type Hub struct {
	mqservice  *queueservice.RabbitMQService
	httpServer *httpservice.HttpServer
	matoptPool sync.Pool
}

func NewRunner() Runner {
	return newHub()
}

func newHub() *Hub {
	h := &Hub{}
	h.matoptPool.New = func() interface{} {
		return opencvservice.NewMatOption()
	}

	return h
}

func (h *Hub) RegisterServices() {
	manage.RegisterService(&queueservice.RabbitMQService{})
	manage.RegisterService(&httpservice.HttpServer{})
}

func (h *Hub) Run() error {
	log.Println("running hub")

	if err := h.services(); err != nil {
		return err
	}

	h.httpServer.ResetRouter(h)

	sub, err := h.mqservice.Subscribe(context.Background())
	if err != nil {

		return err
	}

	var rec *opencvservice.VideoRecord

	go func() {
		for msg := range sub {
			matopt := h.matoptPool.Get().(*opencvservice.MatOption)

			mat, err := matopt.Mat(msg.Body)
			if err != nil {
				log.Fatal(err)

				h.matoptPool.Put(matopt)
				break
			}

			if rec == nil || rec.Size() >= config.ToBytes(10) {
				rec, err = opencvservice.OpenVideoRecord(mat)
				if err != nil {
					log.Fatal(err)
					h.matoptPool.Put(matopt)

					break
				}
				h.matoptPool.Put(matopt)
				continue
			}

			if err := rec.Write(mat); err != nil {
				h.matoptPool.Put(matopt)
				continue
			}
			h.matoptPool.Put(matopt)
		}
		rec.Close()
	}()

	return nil
}

func (h *Hub) Close() error {
	return errors.New("not yet implemented")
}

func (h *Hub) StreamVideo() <-chan httpservice.VideoResponse {
	vrchan := make(chan httpservice.VideoResponse)

	go func() {
		sub, err := h.mqservice.Subscribe(context.Background())
		if err != nil {
			vrchan <- httpservice.VideoResponse{Err: err}
			close(vrchan)
		}

		for msg := range sub {
			matopt := h.matoptPool.Get().(opencvservice.MatOption)

			data, err := matopt.MatCompress(msg.Body, opencvservice.JPEGCompress)
			if err != nil {
				log.Fatal(err)

				h.matoptPool.Put(matopt)
				break
			}

			v := httpservice.VideoResponse{
				Data: data,
			}

			h.matoptPool.Put(matopt)
			vrchan <- v
		}
		close(vrchan)
	}()

	return vrchan
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

	httpServer, err := manage.GetService("service.http.server")
	if err != nil {
		return fmt.Errorf("service.http.server instance error: %s", err.Error())
	}

	httpinstance := httpServer.Intstance

	hsrv, ok := httpinstance.(*httpservice.HttpServer)
	if !ok {
		return errors.New("http server service assertion problem")
	}

	h.httpServer = hsrv

	h.mqservice = mq

	return nil
}
