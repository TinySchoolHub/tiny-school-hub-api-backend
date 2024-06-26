package main

import (
	"log"
	"net/http"
	"github.com/TinySchoolHub/tiny-school-hub-api-backend/routers"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/TinySchoolHub/tiny-school-hub-api-backend/docs"
	
)

// @title TinySchoolHub Backend API
// @version 1.0
// @description This is a sample server for managing users.
// @host localhost:8080
// @BasePath /
func main() {
    r := routers.SetupRouter()

	// Swagger endpoint
	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
    log.Fatal(http.ListenAndServe(":8080", r))

}