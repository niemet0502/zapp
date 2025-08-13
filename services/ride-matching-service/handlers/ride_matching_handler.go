package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/niemet0502/zapp/services/ride-matching-service/services"
	"github.com/thedevsaddam/govalidator"

	"github.com/niemet0502/zapp/pkg/connection"
)

type RideMatchingHandler struct {
	svc *services.RideMatchingService
}

func CreateRideMatchingHandler(svc *services.RideMatchingService) *RideMatchingHandler {
	return &RideMatchingHandler{svc}
}

func (h *RideMatchingHandler) RideMatching(w http.ResponseWriter, r *http.Request) {
	rideId := chi.URLParam(r, "rideId")

	i, err := strconv.Atoi(rideId)
	if err != nil {
		log.Fatal(err)
	}

	h.svc.Matching(uint(i))
}

func (h *RideMatchingHandler) RandomMessageSSEHandler(w http.ResponseWriter, r *http.Request) {
	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin")) // Optional for testing

	clientIP := r.RemoteAddr
	log.Printf("New SSE client connected: %s", clientIP)

	driverID := r.URL.Query().Get("id")
	userType := r.URL.Query().Get("type")

	if driverID == "" || userType == "" {
		http.Error(w, "Driver ID required", http.StatusBadRequest)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	ch := make(chan string)

	connection.DriversMu.Lock()
	connection.Drivers[driverID] = connection.UserConnected{ID: driverID, Type: userType, Chan: ch}
	connection.DriversMu.Unlock()

	log.Printf("Driver %s connected", driverID)

	ctx := r.Context()
	go func() {
		<-ctx.Done()
		connection.DriversMu.Lock()
		delete(connection.Drivers, driverID)
		close(ch)
		connection.DriversMu.Unlock()
		log.Printf("Driver %s disconnected", driverID)
	}()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-ch:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		case <-ticker.C:
			fmt.Fprintf(w, ": keep-alive\n\n")
			flusher.Flush()
		}
	}

}

func (h *RideMatchingHandler) AcceptRide(w http.ResponseWriter, r *http.Request) {
	var request services.AcceptRideRequest
	rideId := chi.URLParam(r, "rideId")

	rules := govalidator.MapData{
		"driver_id": []string{"required"},
		"response":  []string{"in:yes,no"},
	}

	opts := govalidator.Options{
		Request: r,
		Data:    &request,
		Rules:   rules,
	}

	v := govalidator.New(opts)
	validationErrors := v.ValidateJSON()

	if len(validationErrors) > 0 {
		slog.Error("Validation failed", slog.Any("errors", validationErrors))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(validationErrors)
		return
	}

	id, _ := strconv.Atoi(rideId)

	result, err := h.svc.AcceptRide(uint(id), request)

	if err != nil {
		slog.Error(err.Error())
		return
	}

	res, _ := json.Marshal(result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (h *RideMatchingHandler) RideUpdate(w http.ResponseWriter, r *http.Request) {
	var ride services.RideUpdateRequest

	rules := govalidator.MapData{
		"status":  []string{"required"},
		"ride_id": []string{"required"},
	}

	opts := govalidator.Options{
		Request: r,
		Data:    &ride,
		Rules:   rules,
	}

	v := govalidator.New(opts)
	validationErrors := v.ValidateJSON()

	if len(validationErrors) > 0 {
		slog.Error("Validation failed", slog.Any("errors", validationErrors))
		http.Error(w, "Validation failed", http.StatusBadRequest)
		return
	}

	result, err := h.svc.RideUpdate(ride)

	if err != nil {
		slog.Error(err.Error())
		return
	}

	res, _ := json.Marshal(result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
