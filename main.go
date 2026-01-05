package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	store := NewServiceStore()
	notifier := NewTelegramNotifier(
		"BOT_TOKEN",
		"CHAT_ID",
	)
	manager := NewServiceManager(
		store,
		notifier,
		1*time.Minute,
	)

	go manager.Start()
	notifier.NotifyDown(&Service{Name: "FastAPI1"})
	notifier.NotifyUp(&Service{Name: "FastAPI1"})

	http.HandleFunc(
		"/services",
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			switch r.Method {
			case http.MethodGet:
				ListServicesHandler(w, r, store)
			case http.MethodPost:
				AddServiceHandler(w, r, store)
			case http.MethodDelete:
				DeleteServiceHandler(w, r, store)
			default:
				http.Error(
					w,
					"Method not allowed",
					http.StatusMethodNotAllowed,
				)
			}
		})
	fmt.Println("Healthcheck service started on :8081")
	log.Fatal(
		http.ListenAndServe(
			"127.0.0.1:8081",
			nil,
		),
	)
}

func ListServicesHandler(
	w http.ResponseWriter,
	r *http.Request,
	store *ServiceStore,
) {
	service := store.GetAll()
	w.Header().Set(
		"Content-Type",
		"application/json",
	)
	json.NewEncoder(w).Encode(service)
	LogEvent("Listed all services")
}

func AddServiceHandler(
	w http.ResponseWriter,
	r *http.Request,
	store *ServiceStore) {
	var service Service
	if err := json.NewDecoder(
		r.Body,
	).Decode(
		&service,
	); err != nil {
		http.Error(
			w,
			"Invalid request",
			http.StatusBadRequest,
		)
		LogEvent("Failed to add service: invalid request")
		return
	}
	store.Add(&service)
	LogEvent("Added service: " + service.Name + " (" + service.URL + ")")
	w.WriteHeader(http.StatusCreated)
}

func DeleteServiceHandler(
	w http.ResponseWriter,
	r *http.Request,
	store *ServiceStore,
) {
	var service Service
	if err := json.NewDecoder(
		r.Body,
	).Decode(
		&service,
	); err != nil {
		http.Error(
			w,
			"Invalid request",
			http.StatusBadRequest,
		)
		LogEvent("Failed to delete service: invalid request")
		return
	}
	store.Delete(service.Name)
	LogEvent("Deleted service: " + service.Name)
	w.WriteHeader(http.StatusNoContent)
}
