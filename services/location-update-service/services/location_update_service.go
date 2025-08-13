package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"

	"github.com/niemet0502/zapp/pkg/connection"
	"github.com/niemet0502/zapp/pkg/events"
	"github.com/redis/go-redis/v9"
)

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
type Update struct {
	Type      string   `json:"type"`
	DriverId  string   `json:"driver_id"`
	Timestamp string   `json:"timestamp"`
	Location  Location `json:"location"`
}
type LocationUpdateService struct {
	redisClient *redis.Client
}

func CreateLocationUpdateService(redisClient *redis.Client) *LocationUpdateService {
	return &LocationUpdateService{redisClient}
}

func (s *LocationUpdateService) UpdateLocation(update Update) (string, error) {
	ctx := context.Background()

	_, err := s.redisClient.GeoAdd(ctx, "drivers:location", &redis.GeoLocation{
		Longitude: update.Location.Longitude,
		Latitude:  update.Location.Latitude,
		Name:      update.DriverId,
	}).Result()

	if err != nil {
		return "", nil
	}

	// publish a message so the ride_matching_service send the update to the connected customer
	connection.ActiveRideMu.Lock()
	data, ok := connection.ActiveRide[update.DriverId]
	connection.ActiveRideMu.Unlock()

	if ok {
		key := fmt.Sprintf("driver:%s:location", update.DriverId)
		var obj events.DriverLocationEvent

		obj.Lat = float32(update.Location.Latitude)
		obj.Long = float32(update.Location.Longitude)
		obj.DriverId = update.DriverId
		obj.RiderId = data

		result, _ := json.Marshal(obj)

		err := s.redisClient.Publish(ctx, key, result).Err()

		if err != nil {
			slog.Info("failed to send the driver location update")
		}
	}

	return "Location successfully updated", nil
}

func (s *LocationUpdateService) SubscribeToRideChannel() {
	ctx := context.Background()

	pubsub := s.redisClient.Subscribe(ctx, "ride:started")

	go func() {
		for msg := range pubsub.Channel() {
			var event events.RideStartedMessage

			if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
				log.Printf("Failed to parse ride:started payload: %v", err)
				continue
			}

			log.Printf("Received ride:started for driver %s and rider %s",
				event.DriverId, event.RideId)

			// add the driverId into the map
			connection.ActiveRideMu.Lock()
			connection.ActiveRide[event.DriverId] = event.RideId
			connection.ActiveRideMu.Unlock()
		}
	}()
}

func (s *LocationUpdateService) SubscribeToRideUpdateChannel() {
	ctx := context.Background()

	pubsub := s.redisClient.Subscribe(ctx, "ride:completed")

	go func() {
		for msg := range pubsub.Channel() {
			var event events.RideStartedMessage

			if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
				log.Printf("Failed to parse ride:started payload: %v", err)
				continue
			}

			log.Printf("Received ride:completed for driver %s and rider %s",
				event.DriverId, event.RideId)

			// add the driverId into the map
			connection.ActiveRideMu.Lock()
			_, ok := connection.ActiveRide[event.DriverId]

			if ok {
				delete(connection.ActiveRide, event.DriverId)
			}
			connection.ActiveRideMu.Unlock()
		}
	}()
}
