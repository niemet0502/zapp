package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RideStatus string

const (
	StatusRequested         RideStatus = "REQUESTED"
	StatusDriverAssigned    RideStatus = "DRIVER_ASSIGNED"
	StatusDriverEnRoute     RideStatus = "DRIVER_EN_ROUTE"
	StatusArrived           RideStatus = "ARRIVED"
	StatusInProgress        RideStatus = "IN_PROGRESS"
	StatusCompleted         RideStatus = "COMPLETED"
	StatusCancelledByRider  RideStatus = "CANCELLED_BY_RIDER"
	StatusCancelledByDriver RideStatus = "CANCELLED_BY_DRIVER"
	StatusNoShow            RideStatus = "NO_SHOW"
	StatusFailed            RideStatus = "FAILED"
)

type Ride struct {
	ID                   uint           `gorm:"primaryKey;autoIncrement;index"`
	PickupLocationLat    float64        `gorm:"column:pickup_location_lat;;not null" json:"pickup_location_lat"`
	PickupLocationLong   float64        `gorm:"column:pickup_location_long;;not null" json:"pickup_location_long"`
	DestinationLat       float64        `gorm:"column:destination_lat;not null" json:"destination_lat"`
	DestinationLong      float64        `gorm:"column:destination_long;not null" json:"destination_long"`
	DriverId             *uuid.UUID     `gorm:"column:driver_id" json:"driver_id"`
	RiderId              uuid.UUID      `gorm:"column:rider_id;not null" json:"rider_id"`
	Status               RideStatus     `gorm:"column:status;default:REQUESTED; not null" json:"status"`
	EstimatedPickupTime  string         `gorm:"column:estimated_pickup_time;type:varchar(100);not null" json:"estimated_pickup_time"`
	EstimatedDropoffTime string         `gorm:"column:estimated_dropoff_time;type:varchar(100);not null" json:"estimated_dropoff_time"`
	CreatedAt            time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt            time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt            gorm.DeletedAt `gorm:"index;column:deleted_at"`
}
