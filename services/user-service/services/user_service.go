package services

import (
	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/niemet0502/zapp/pkg/models"
)

type UserService struct {
	db *gorm.DB
}

type UserRequest struct {
	Email    string  `json:"email"`
	FullName string  `json:"fullName"`
	UserType string  `json:"userType"`
	CarInfo  *string `json:"carInfo,omitempty"`
	Status   *string `json:"status,omitempty"`
}

func CreateUserService(db *gorm.DB) *UserService {
	return &UserService{db}
}

func (svc *UserService) CreateUser(userRequest UserRequest) (*models.User, error) {
	user := models.User{ID: uuid.New(), Email: userRequest.Email, FullName: userRequest.FullName, UserType: userRequest.UserType, CarInfo: userRequest.CarInfo, Status: userRequest.Status}
	result := svc.db.Create(&user)

	if result.Error != nil {
		return &user, result.Error
	}

	return &user, nil
}
