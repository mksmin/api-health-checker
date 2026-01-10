package storage

import (
	"encoding/json"
	"fmt"
	"healthchecker/internal/common"
	"healthchecker/internal/logs"
	"os"
	"path/filepath"
	"sync"
)

type ServiceRepository interface {
	Load() (map[string]*common.Service, error)
	Save(map[string]*common.Service) error
}

type ServiceStore struct {
	services map[string]*common.Service
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
	service *common.Service,
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
) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.services[name]; !exists {
		return false
	}

	delete(s.services, name)
	s.repo.Save(s.services)
	return true
}

func (
	s *ServiceStore,
) GetAll() []*common.Service {
	s.mu.RLock()
	defer s.mu.RUnlock()

	all := make(
		[]*common.Service,
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
	map[string]*common.Service,
	error,
) {
	if _, err := os.Stat(s.path); os.IsNotExist(err) {
		return map[string]*common.Service{}, nil
	}

	data, err := os.ReadFile(s.path)
	if err != nil {
		return nil, err
	}

	var services map[string]*common.Service
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
	services map[string]*common.Service,
) error {
	dir := filepath.Dir(s.path)
	tmp := s.path + ".tmp"

	data, err := json.MarshalIndent(
		services,
		"",
		" ",
	)
	if err != nil {
		logs.LogEvent(
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
		logs.LogEvent(
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
		logs.LogEvent(
			fmt.Sprintf(
				"Failed to WriteFile: %s", err,
			),
		)
		return err
	}

	logs.LogEvent(fmt.Sprintf("Saving services to %s", tmp))

	return os.Rename(
		tmp,
		s.path,
	)
}
