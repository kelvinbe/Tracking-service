package main

import (
	"context"
	"log"
	"tracking-service/repository"
	"tracking-service/utils"
	"time"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/getsentry/sentry-go"
)

func main() {
	utils.LoadEnv()

	err := sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
		Environment: os.Getenv("SENTRY_ENV"),
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)

	sentry.CaptureMessage("It works for the second time!") 

	// this will initialize the repository and load the database connection into it
	repo := repository.InitRepository()
	
	// close the connection when main exits
	defer repo.Client.Disconnect(context.Background())

	// create a new fiber app
	app := fiber.New()

	//setting up the routes for the fiber app
	repo.SetupRotes(app)

	app.Listen("0.0.0.0:8080")

	log.Printf("App listening on port %d", 8080);

}