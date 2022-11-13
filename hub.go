//go:build hub

package easymotion

import (
	"errors"

	"github.com/rlaskowski/easymotion/service/queueservice"
	"github.com/rlaskowski/manage"
)

type Hub struct {
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
	return errors.New("not yet implemented")
}

func (h *Hub) Close() error {
	return errors.New("not yet implemented")
}
