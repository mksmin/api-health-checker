package main

import (
	"encoding/json"
	"healthchecker/internal/common"
	"healthchecker/internal/logs"
	"healthchecker/internal/notifier"
	"healthchecker/internal/services"
	"healthchecker/internal/storage"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("Warning: could not load .env file: %v", err)
	}

	duration := 60 * time.Second
	if envDuration := os.Getenv("SERVICES_DURATION"); envDuration != "" {
		if parseEnvDuration, err := time.ParseDuration(envDuration); err != nil {
			duration = parseEnvDuration
		}
	}

	path := os.Getenv("SERVICES_FILE")
	if path == "" {
		path = "./data/services.json"
	}

	repo := storage.NewJSONStore(path)
	store, err := storage.NewServiceStore(repo)
	notify := notifier.NewTelegramNotifier()
	manager := services.NewServiceManager(
		store,
		notify,
		duration,
	)

	if err != nil {
		log.Fatal(err)
	}

	go manager.Start()
	notify.NotifyUp(&common.Service{Name: "Health Checker"})

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
	logs.LogEvent("Healthcheck service started")

	addr := os.Getenv("HTTP_ADDR")

	if addr == "" {
		addr = ":8081"
	}

	log.Fatal(
		http.ListenAndServe(
			addr,
			nil,
		),
	)
}

func ListServicesHandler(
	w http.ResponseWriter,
	r *http.Request,
	store *storage.ServiceStore,
) {
	service := store.GetAll()
	w.Header().Set(
		"Content-Type",
		"application/json",
	)
	json.NewEncoder(w).Encode(service)
	logs.LogEvent("Listed all services")
}

func AddServiceHandler(
	w http.ResponseWriter,
	r *http.Request,
	store *storage.ServiceStore,
) {
	var service common.Service
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
		logs.LogEvent("Failed to add service: invalid request")
		return
	}
	store.Add(&service)
	logs.LogEvent("Added service: " + service.Name + " (" + service.URL + ")")
	w.WriteHeader(http.StatusCreated)
}

func DeleteServiceHandler(
	w http.ResponseWriter,
	r *http.Request,
	store *storage.ServiceStore,
) {
	var service common.Service
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
		logs.LogEvent("Failed to delete service: invalid request")
		return
	}
	result := store.Delete(service.Name)
	if result {
		logs.LogEvent("Deleted service: " + service.Name)
		w.WriteHeader(http.StatusNoContent)
	} else {
		logs.LogEvent("Delete: Service not found: " + service.Name)
		w.WriteHeader(http.StatusNotFound)
	}
}
