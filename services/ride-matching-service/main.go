package main

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/niemet0502/zapp/services/ride-matching-service/handlers"
	"github.com/niemet0502/zapp/services/ride-matching-service/services"
)

func main() {
	r := chi.NewRouter()

	slog.Info("ride matching service")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // No password set
		DB:       0,  // Use default DB
		Protocol: 2,
	})

	dsn := "host=localhost user=user password=password dbname=mydatabase port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	svc := services.CreateRideMatchingService(redisClient, db)

	svc.SubscribeToDriverLocationUpdate()

	h := handlers.CreateRideMatchingHandler(svc)

	r.Patch("/ride/request/{rideId}", h.RideMatching)

	r.Patch("/ride/accept/{rideId}", h.AcceptRide)

	r.Patch("/ride/update/{rideId}", h.RideUpdate)

	r.HandleFunc("/sse", h.RandomMessageSSEHandler)

	if err := http.ListenAndServe(":3004", r); err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
