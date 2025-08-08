package services

import (
	"github.com/google/uuid"
	"github.com/niemet0502/zapp/pkg/models"
	"gorm.io/gorm"
)

type RideService struct {
	db *gorm.DB
}

type RideRequest struct {
	PickupLocationLat    float64   `json:"pickup_location_lat"`
	PickupLocationLong   float64   `json:"pickup_location_long"`
	DestinationLat       float64   `json:"destination_lat"`
	DestinationLong      float64   `json:"destination_long"`
	RideId               uuid.UUID `json:"rider_id"`
	DriverId             uint32    `json:"driver_id"`
	EstimatedPickupTime  string    `json:"estimated_pickup_time"`
	EstimatedDropoffTime string    `json:"estimated_dropoff_time"`
}

type AcceptRideRequest struct {
	RideId   uuid.UUID `json:"ride_id"`
	DriverId uuid.UUID `json:"driver_id"`
	Response string    `json:"response"`
}

func CreateRideService(db *gorm.DB) *RideService {
	return &RideService{db}
}
func (s *RideService) CreateRide(request RideRequest) (*models.Ride, error) {
	ride := models.Ride{PickupLocationLat: request.PickupLocationLat, PickupLocationLong: request.PickupLocationLong, DestinationLat: request.DestinationLat, DestinationLong: request.DestinationLong, RiderId: request.RideId}

	result := s.db.Create(&ride)

	if result.Error != nil {
		return &ride, result.Error
	}

	return &ride, nil
}

func (s *RideService) AcceptRide(request AcceptRideRequest) (*models.Ride, error) {
	var ride models.Ride
	s.db.First(&ride, request.RideId)

	ride.DriverId = &request.DriverId

	s.db.Save(&ride)

	// send the rider a notification

	//
	return &ride, nil
}
