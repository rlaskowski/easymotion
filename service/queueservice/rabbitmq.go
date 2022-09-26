package queueservice

import "github.com/rlaskowski/manage"

type RabbitMQService struct {
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
	return nil
}

func (r *RabbitMQService) Stop() error {
	return nil
}
