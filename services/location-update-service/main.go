package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/niemet0502/zapp/services/location-update-service/controllers"
	"github.com/niemet0502/zapp/services/location-update-service/services"
	"github.com/redis/go-redis/v9"
)

func main() {
	slog.Info("location update service")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // No password set
		DB:       0,  // Use default DB
		Protocol: 2,  // Connection protocol
	})

	locationService := services.CreateLocationUpdateService(redisClient)
	locationHandler := controllers.CreateLocationUpdatHandler(locationService)

	locationService.SubscribeToRideChannel()
	locationService.SubscribeToRideUpdateChannel()

	http.HandleFunc("/ws", locationHandler.WsHandler)

	fmt.Println("Server listening on :3002")
	if err := http.ListenAndServe(":3002", nil); err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
