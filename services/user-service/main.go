package main

import (
	"net/http"

	"github.com/niemet0502/zapp/services/user-service/routes"
)

func main() {
	println("ride sharing app ")

	r := routes.ApiServer()

	http.ListenAndServe(":3000", r)
}
