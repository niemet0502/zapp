package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserType string

const (
	RiderType  UserType = "rider"
	DriverType UserType = "driver"
)

// Enum for Driver Status
type DriverStatus string

const (
	Available DriverStatus = "available"
	Busy      DriverStatus = "busy"
	Offline   DriverStatus = "offline"
	Suspended DriverStatus = "suspended"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Email     string         `gorm:"unique;not null" json:"email"`
	FullName  string         `gorm:"not null" json:"fullName"`
	UserType  string         `gorm:":not null" json:"userType"`
	CarInfo   *string        `gorm:"type:text" json:"carInfo,omitempty"` // Nullable for Riders, omit if nil
	Status    *string        `gorm:"not null" json:"status,omitempty"`   // Nullable for Riders, omit if nil
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at"`
}
