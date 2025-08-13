package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/niemet0502/zapp/pkg/events"
	"github.com/niemet0502/zapp/pkg/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/niemet0502/zapp/pkg/connection"
)

type AcceptRideRequest struct {
	DriverId uuid.UUID `json:"driver_id"`
	Response string    `json:"response"`
}

type RideUpdateRequest struct {
	RideId uint              `json:"ride_id"`
	Status models.RideStatus `json:"status"`
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
	var ride models.Ride

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

		var rideState models.Ride

		s.db.Find(&rideState, rideId)

		if rideState.Status == models.StatusDriverEnRoute {
			return
		}

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
	ctx := context.Background()
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

	// publish the message ride started with the driverid and the riderId
	var payload events.RideStartedMessage
	payload.DriverId = request.DriverId.String()
	payload.RideId = ride.RiderId.String()

	result, _ := json.Marshal(payload)

	err := s.redisClient.Publish(ctx, "ride:started", result).Err()

	if err != nil {
		slog.Info("failed to publish the ride:started event")
	}

	return &ride, nil
}

func (s *RideMatchingService) SubscribeToDriverLocationUpdate() {
	ctx := context.Background()

	pubsub := s.redisClient.PSubscribe(ctx, "driver:*:location")

	go func() {
		for msg := range pubsub.Channel() {
			var event events.DriverLocationEvent

			if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
				log.Printf("Failed to parse ride:started payload: %v", err)
				continue
			}

			log.Printf("Driver %s is moving %s",
				event.DriverId, event.RiderId)

			// check if the rider exist in the local map
			// if so send in the channel the driver location update
			connection.DriversMu.Lock()
			channel, ok := connection.Drivers[event.RiderId]
			connection.DriversMu.Unlock()

			if ok {
				message := fmt.Sprintf("The driver %s - lat %f long %f", event.DriverId, event.Lat, event.Long)
				channel.Chan <- message
			}
		}
	}()
}

func (s *RideMatchingService) RideUpdate(rideUpdateRequest RideUpdateRequest) (*models.Ride, error) {
	ctx := context.Background()
	var ride models.Ride
	s.db.First(&ride, rideUpdateRequest.RideId)

	println(rideUpdateRequest.Status)

	if rideUpdateRequest.Status == models.StatusArrived {
		connection.DriversMu.Lock()
		rider, ok := connection.Drivers[ride.RiderId.String()]
		connection.DriversMu.Unlock()

		if ok {
			rider.Chan <- "Your driver is here"
		}
	}

	if rideUpdateRequest.Status == models.StatusCompleted {
		var payload events.RideStartedMessage

		payload.DriverId = ride.DriverId.String()
		payload.RideId = ride.RiderId.String()

		result, _ := json.Marshal(payload)

		err := s.redisClient.Publish(ctx, "ride:completed", result).Err()

		if err != nil {
			slog.Info("failed to publish the event ride completed")
		}
	}

	ride.Status = rideUpdateRequest.Status
	s.db.Save(&ride)

	return &ride, nil
}
