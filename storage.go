package main

import "sync"

type ServiceStore struct {
	services map[string]*Service
	mu       sync.RWMutex
}

func NewServiceStore() *ServiceStore {
	return &ServiceStore{
		services: make(
			map[string]*Service,
		),
	}
}

func (
	s *ServiceStore,
) Add(
	service *Service,
) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.services[service.Name] = service
}

func (
	s *ServiceStore,
) Delete(
	name string,
) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.services, name)
}

func (
	s *ServiceStore,
) GetAll() []*Service {
	s.mu.RLock()
	defer s.mu.RUnlock()

	all := make(
		[]*Service,
		0,
		len(s.services),
	)

	for _, service := range s.services {
		all = append(all, service)
	}
	return all
}
