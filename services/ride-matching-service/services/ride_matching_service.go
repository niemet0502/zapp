package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/niemet0502/zapp/pkg/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/niemet0502/zapp/pkg/connection"
)

type AcceptRideRequest struct {
	DriverId uuid.UUID `json:"driver_id"`
	Response string    `json:"response"`
}

type RideMatchingService struct {
	redisClient *redis.Client
	db          *gorm.DB
}

func CreateRideMatchingService(redisClient *redis.Client,
	db *gorm.DB) *RideMatchingService {
	return &RideMatchingService{redisClient, db}
}

func (s *RideMatchingService) Matching(rideId uint) {
	slog.Info("called matching")
	var ride models.Ride

	println("rideid %d", rideId)
	// fetch ride from the db
	s.db.First(&ride, rideId)

	ctx := context.Background()
	// fetch nearby drivers
	drivers, err := s.redisClient.GeoSearch(ctx, "drivers:location", &redis.GeoSearchQuery{
		Longitude:  ride.PickupLocationLong,
		Latitude:   ride.PickupLocationLat,
		Radius:     5,
		RadiusUnit: "km",
	}).Result()

	if err != nil {
		slog.Error(err.Error())
		slog.Error("Failed to fetch drivers")
		return
	}

	for _, driverID := range drivers {

		connection.DriversMu.Lock()
		driver, ok := connection.Drivers[driverID]

		if !ok {
			println("driver %s isn't connected", driverID)
		}

		key := fmt.Sprintf("drivers:%s", driverID)

		// check if the driver is locked if so move to the next one
		ttl, err := s.redisClient.TTL(ctx, key).Result()

		if err != nil {
			slog.Info("failed to check driver in the TTL")
		}

		if ttl.Seconds() < 0 {
			result, err := s.redisClient.Set(ctx, key, "notified", 20*time.Second).Result()

			println("result %s", result)
			if err != nil {
				slog.Error("Failed to add the driver in the redis TTL")
			}
			// send notification to driver using SSE
			driver.Chan <- fmt.Sprintf("Here is a new ride for you %d", ride.ID)
			connection.DriversMu.Unlock()

			time.Sleep(20 * time.Second)
		}
	}

}

func (s *RideMatchingService) AcceptRide(rideId uint, request AcceptRideRequest) (*models.Ride, error) {
	slog.Info("called")
	var ride models.Ride
	s.db.First(&ride, rideId)

	ride.DriverId = &request.DriverId
	ride.Status = models.StatusDriverEnRoute

	s.db.Save(&ride)
	slog.Info(ride.RiderId.String())
	// send the rider a notification

	connection.DriversMu.Lock()
	customer, ok := connection.Drivers[ride.RiderId.String()]

	connection.DriversMu.Unlock()

	if !ok {
		slog.Info("The rider isn't connected")
		return &ride, fmt.Errorf("the rider isn't connected")
	}

	customer.Chan <- "The driver is coming !!!"

	// update the ride status
	return &ride, nil
}
