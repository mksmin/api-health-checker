package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type ServiceRepository interface {
	Load() (map[string]*Service, error)
	Save(map[string]*Service) error
}

type ServiceStore struct {
	services map[string]*Service
	mu       sync.RWMutex
	repo     ServiceRepository
}

func NewServiceStore(
	repo ServiceRepository,
) (*ServiceStore, error) {
	services, err := repo.Load()
	if err != nil {
		return nil, err
	}

	return &ServiceStore{
		services: services,
		repo:     repo,
	}, nil
}

func (
	s *ServiceStore,
) Add(
	service *Service,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.services[service.Name] = service
	return s.repo.Save(s.services)
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

type JSONStore struct {
	path string
}

func NewJSONStore(path string) *JSONStore {
	return &JSONStore{path: path}
}

func (s *JSONStore) Load() (
	map[string]*Service,
	error,
) {
	if _, err := os.Stat(s.path); os.IsNotExist(err) {
		return map[string]*Service{}, nil
	}

	data, err := os.ReadFile(s.path)
	if err != nil {
		return nil, err
	}

	var services map[string]*Service
	if err := json.Unmarshal(
		data,
		&services,
	); err != nil {
		return nil, err
	}

	return services, nil
}

func (
	s *JSONStore,
) Save(
	services map[string]*Service,
) error {
	dir := filepath.Dir(s.path)
	tmp := s.path + ".tmp"

	data, err := json.MarshalIndent(
		services,
		"",
		" ",
	)
	if err != nil {
		LogEvent(
			fmt.Sprintf(
				"Failed to MarshalIndent: %s", err,
			),
		)
		return err
	}

	if err := os.MkdirAll(
		dir,
		0777,
	); err != nil {
		LogEvent(
			fmt.Sprintf(
				"Failed to MkdirAll: %s", err,
			),
		)
		return err
	}

	if err := os.WriteFile(
		tmp,
		data,
		0666,
	); err != nil {
		LogEvent(
			fmt.Sprintf(
				"Failed to WriteFile: %s", err,
			),
		)
		return err
	}

	LogEvent(fmt.Sprintf("Saving services to %s", tmp))

	return os.Rename(
		tmp,
		s.path,
	)
}
