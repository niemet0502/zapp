package controllers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/niemet0502/zapp/services/ride-service/services"
	"github.com/thedevsaddam/govalidator"
)

type RideController struct {
	svc *services.RideService
}

func CreateRideController(svc *services.RideService) *RideController {
	return &RideController{svc}
}

func (c *RideController) RideEstimation(w http.ResponseWriter, r *http.Request) {
	var ride services.RideRequest

	rules := govalidator.MapData{
		"pickup_location_lat":    []string{"required", "float"},
		"pickup_location_long":   []string{"required", "float"},
		"destination_lat":        []string{"required", "float"},
		"destination_long":       []string{"required", "float"},
		"rider_id":               []string{"required"},
		"estimated_pickup_time":  []string{"required", "between:1,100"},
		"estimated_dropoff_time": []string{"required", "between:1,100"},
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

	result, err := c.svc.CreateRide(ride)
	if err != nil {
		slog.Error("Failed to create ride", slog.Any("error", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	res, _ := json.Marshal(result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func (c *RideController) AcceptRide(w http.ResponseWriter, r *http.Request) {
	var request services.AcceptRideRequest

	rules := govalidator.MapData{
		"ride_id":   []string{"required"},
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

	// Handle validation errors
	if len(validationErrors) > 0 {
		slog.Error("Validation failed", slog.Any("errors", validationErrors))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(validationErrors)
		return
	}

	result, err := c.svc.AcceptRide(request)

	if err != nil {
		slog.Error(err.Error())
		return
	}

	res, _ := json.Marshal(result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
