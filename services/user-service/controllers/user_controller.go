package controllers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/niemet0502/zapp/services/user-service/services"
	"github.com/thedevsaddam/govalidator"
)

type UserController struct {
	svc *services.UserService
}

func CreateUserController(svc *services.UserService) *UserController {
	return &UserController{svc}
}

func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user services.UserRequest

	// Validation rules
	rules := govalidator.MapData{
		"email":    []string{"required", "email"},
		"fullName": []string{"required", "between:3,50"},
		"userType": []string{"required", "in:rider,driver"},
		"carInfo":  []string{"max:100"},
		"status":   []string{"in:available,busy,offline,suspended"},
	}

	opts := govalidator.Options{
		Request: r,
		Data:    &user,
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

	result, err := c.svc.CreateUser(user)

	if err != nil {
		slog.Error("Failed to create user ")
		return
	}

	res, _ := json.Marshal(result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}
