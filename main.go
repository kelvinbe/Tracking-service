package main

import (
	"log"
	"tracking-service/repository"
	"tracking-service/utils"

	"github.com/gofiber/fiber/v2"
)

func main() {
	utils.LoadEnv()

	// this will initialize the repository and load the database connection into it
	repo := repository.InitRepository()
	// close the connection when main exits
	defer repo.Client.Disconnect(*repo.Context)

	// create a new fiber app
	app := fiber.New()

	//setting up the routes for the fiber app
	repo.SetupRotes(app)

	app.Listen("0.0.0.0:8080")

	log.Printf("App listening on port %d", 8080);

}