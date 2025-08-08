package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/niemet0502/zapp/pkg/models"
)

var Connection *gorm.DB

func InitDb() (*gorm.DB, error) {

	dsn := "host=localhost user=user password=password dbname=mydatabase port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}

	db.AutoMigrate(&models.User{})

	Connection = db

	return db, nil
}
