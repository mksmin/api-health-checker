package main

import (
	"net/http"
	"time"
)

type Service struct {
	Name     string
	URL      string
	IsUp     bool
	LastDown time.Time
}

type ServiceManager struct {
	store    *ServiceStore
	notifier *TelegramNotifier
	interval time.Duration
}

func NewServiceManager(
	store *ServiceStore,
	notifier *TelegramNotifier,
	interval time.Duration,
) *ServiceManager {
	return &ServiceManager{
		store:    store,
		notifier: notifier,
		interval: interval,
	}
}

func (
	m *ServiceManager,
) Start() {
	ticker := time.NewTicker(m.interval)
	for range ticker.C {
		services := m.store.GetAll()
		if len(services) == 0 {
			LogEvent("No services to check")
			continue
		}
		for _, s := range services {
			go m.checkService(s)
		}
	}
}

func (
	m *ServiceManager,
) checkService(
	s *Service,
) {
	resp, err := http.Get(s.URL)
	up := err == nil && resp.StatusCode < 500

	if !up && s.IsUp {
		s.LastDown = time.Now()
		m.notifier.NotifyDown(s)
		LogEvent(s.Name + " is down")
	} else if up && !s.IsUp {
		m.notifier.NotifyUp(s)
		LogEvent(s.Name + " is up")
	}

	s.IsUp = up
}
