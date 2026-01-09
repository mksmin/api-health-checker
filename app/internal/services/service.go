package services

import (
	"healthchecker/internal/common"
	"healthchecker/internal/logs"
	"healthchecker/internal/notifier"
	"healthchecker/internal/storage"
	"net/http"
	"time"
)

type ServiceManager struct {
	store    *storage.ServiceStore
	notifier *notifier.TelegramNotifier
	interval time.Duration
}

func NewServiceManager(
	store *storage.ServiceStore,
	notifier *notifier.TelegramNotifier,
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
	defer ticker.Stop()

	for range ticker.C {
		services := m.store.GetAll()
		if len(services) == 0 {
			logs.LogEvent("No services to check")
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
	s *common.Service,
) {
	resp, err := http.Get(s.URL)
	up := err == nil && resp.StatusCode < 500

	if !up && s.IsUp {
		s.LastDown = time.Now()
		m.notifier.NotifyDown(s)
		logs.LogEvent(s.Name + " is down")
	} else if up && !s.IsUp {
		m.notifier.NotifyUp(s)
		logs.LogEvent(s.Name + " is up")
	}

	s.IsUp = up
}
