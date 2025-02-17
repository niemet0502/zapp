package main

import (
	"log/slog"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/niemet0502/zapp/services/user-service/db"
	"github.com/niemet0502/zapp/services/user-service/routes"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		slog.Info("Failed to fetch the env variable")
	}
	db.InitDb()

	r := routes.ApiServer()

	slog.Info("The service is listening on 3000 port")
	http.ListenAndServe(":3000", r)
}
