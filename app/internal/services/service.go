package services

import (
	"fmt"
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
	logs.LogEvent(
		fmt.Sprintf(
			"Checking service: %s (URL: %s)",
			s.Name,
			s.URL,
		),
	)
	resp, err := http.Get(s.URL)
	up := err == nil && resp.StatusCode < 500

	if !up && s.IsUp {
		s.LastDown = time.Now()
		m.notifier.NotifyDown(s)
		logs.LogEvent(
			fmt.Sprintf(
				"Service is down: %s (URL: %s, Status: %v, Error: %v)",
				s.Name, s.URL, getResponseStatus(resp), err,
			),
		)
	} else if up && !s.IsUp {
		m.notifier.NotifyUp(s)
		logs.LogEvent(
			fmt.Sprintf(
				"Service is up: %s (URL: %s, Status: %v)",
				s.Name, s.URL, getResponseStatus(resp),
			),
		)
	}

	if s.IsUp != up {
		s.IsUp = up
		err = m.store.Add(s)
		if err != nil {
			logs.LogEvent(
				fmt.Sprintf("Failed to save service status: %v", err),
			)
		}
	}
}

func getResponseStatus(
	r *http.Response,
) string {
	if r == nil {
		return "no response"
	}
	return fmt.Sprintf("%d %s", r.StatusCode, r.Status)
}
