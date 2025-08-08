package main

import (
	"log/slog"

	"net/http"

	"github.com/niemet0502/zapp/services/ride-service/controllers"
	"github.com/niemet0502/zapp/services/ride-service/routes"
	"github.com/niemet0502/zapp/services/ride-service/services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/niemet0502/zapp/pkg/models"
)

func main() {
	dsn := "host=localhost user=user password=password dbname=mydatabase port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		slog.Error("db connexion failed")
		panic(err)
	}
	db.AutoMigrate(&models.Ride{})

	slog.Info("ride service ")

	r := routes.ApiServer()

	rideSvc := services.CreateRideService(db)
	rideController := controllers.CreateRideController(rideSvc)

	r.Post("/ride/far-estimation", rideController.RideEstimation)

	http.ListenAndServe(":3001", r)

}
