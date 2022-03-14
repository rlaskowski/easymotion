package easymotion

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

var (
	srvRWMutex = &sync.RWMutex{}
	serviceMap = make(map[string]*ServiceInfo)
)

type Service struct {
}

type ServiceRunner interface {
	Start() error
	Stop() error
	CreateService() *ServiceInfo
}

type ServiceInfo struct {
	ID        string
	Intstance ServiceRunner
}

func NewService() *Service {
	return &Service{}
}

func RegisterService(runner ServiceRunner) error {
	service := runner.CreateService()

	if service.ID == "" {
		return errors.New("no service id")
	}

	srvRWMutex.Lock()
	defer srvRWMutex.Unlock()

	if _, ok := serviceMap[service.ID]; ok {
		return fmt.Errorf("service %s has been already registered", service.ID)
	}

	serviceMap[service.ID] = service

	return nil
}

func GetService(id string) (*ServiceInfo, error) {
	srvRWMutex.Lock()
	defer srvRWMutex.Unlock()

	service, ok := serviceMap[id]
	if !ok {
		return nil, fmt.Errorf("service %s not found", id)
	}

	return service, nil
}

func GetServiceID(instance interface{}) string {
	srvRWMutex.Lock()
	defer srvRWMutex.Unlock()

	service, ok := instance.(ServiceInfo)
	if !ok {
		return ""
	}

	return service.ID
}

// Starting all services
func (s *Service) Start() error {
	var result error

	for _, v := range serviceMap {
		if result := v.Intstance.Start(); result != nil {
			log.Printf("couldn't start service %s due to: %s", v.ID, result.Error())
		}
	}

	if result == nil {
		log.Println("All services has been started")
	}

	return result
}

// Stopping all services
func (s *Service) Stop() error {
	var result error

	log.Println("Stopping all services...")

	for _, v := range serviceMap {
		if result := v.Intstance.Stop(); result != nil {
			log.Printf("couldn't start service %s due to: %s", v.ID, result.Error())
		}
	}

	return result
}
