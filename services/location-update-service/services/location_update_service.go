package services

import (
	"context"

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

	res := s.redisClient.GeoAdd(ctx, "drivers:location", &redis.GeoLocation{
		Longitude: update.Location.Longitude,
		Latitude:  update.Location.Latitude,
		Name:      update.DriverId,
	})

	if res == nil {
		return "", nil
	}

	return "Location successfully updated", nil
}
