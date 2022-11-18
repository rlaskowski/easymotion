package queueservice

import (
	"context"
	"log"

	"github.com/rlaskowski/easymotion/config"
	"github.com/rlaskowski/manage"
	"github.com/rlaskowski/manage/rabbitmq"
)

type RabbitMQService struct {
	mqservice *rabbitmq.RabbitMQService
}

func (RabbitMQService) CreateService() *manage.ServiceInfo {
	return &manage.ServiceInfo{
		ID:        "service.rabbitmq",
		Priority:  2,
		Intstance: newRabbitMQ(),
	}
}

func newRabbitMQ() *RabbitMQService {
	return &RabbitMQService{}
}

func (r *RabbitMQService) Start() error {
	log.Println("starting rabbitmq service")

	mqservice, err := rabbitmq.Open(config.DefaultOptions.MQAddress)
	if err != nil {
		return err
	}

	r.mqservice = mqservice

	return nil
}

func (r *RabbitMQService) Stop() error {
	log.Println("stopping rabbitmq service")

	return r.mqservice.Close()
}

func (r *RabbitMQService) Publish(ctx context.Context, msg rabbitmq.Message) error {
	options := rabbitmq.PubOptions{
		Exchange: rabbitmq.Exchange{
			Name: "easymotion.camera.fanout",
			Kind: rabbitmq.FanoutExchange,
		},
		Message: msg,
	}

	return r.mqservice.Publish(ctx, options)
}

func (r *RabbitMQService) Subscribe(ctx context.Context) (<-chan rabbitmq.Message, error) {
	options := rabbitmq.SubOptions{
		Exchange: rabbitmq.Exchange{
			Name: "easymotion.camera.fanout",
			Kind: rabbitmq.FanoutExchange,
		},
	}

	return r.mqservice.Subscribe(ctx, options)
}
