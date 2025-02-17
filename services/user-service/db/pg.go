package db

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Connection *gorm.DB

func InitDb() error {

	port := os.Getenv("POSTGRES_PORT")
	dbname := os.Getenv("POSTGRES_DB")
	password := os.Getenv("POSTGRES_PASSWORD")
	user := os.Getenv("POSTGRES_USER")
	host := os.Getenv("POSTGRES_HOST")

	p, _ := strconv.Atoi(port)

	psqlInfo := fmt.Sprintf("host=%s user=%s  "+
		"password=%s dbname=%s port=%d sslmode=disable",
		host, user, password, dbname, int(p))

	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		slog.Error(err.Error())
	}

	Connection = db

	return nil
}
