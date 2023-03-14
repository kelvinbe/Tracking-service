package repository

import (
	"context"
	"log"
	"tracking-service/handlers"
	"tracking-service/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	Client *mongo.Client
	DB    *mongo.Database
	Context *context.Context
}

func InitRepository() *Repository {

	connection, err := storage.NewClient()

	if err != nil {
		log.Fatalf("Error connecting to database: %s", err.Error())
	}

	err = connection.Client.Connect(*connection.Context)

	if err != nil {
		log.Fatalf("Error connecting to database: %s", err.Error())
	}
	
	db := connection.Client.Database("tracking_service_test")

	return &Repository{
		Client: connection.Client,
		DB:    db,
		Context: connection.Context,
	}
}

// function to make it easier to add routes
func GenerateHandlers( handlers *[]handlers.APIHandler, api *fiber.Router, db *mongo.Database) {
	
	for i := range (*handlers) {
		handler := (*handlers)[i]
		_h := func (ctx *fiber.Ctx) error {
			return handler.Handler(db, ctx)
		}
		switch handler.Method {
		case	"GET":
			(*api).Get(handler.Route, _h)
			log.Printf("Request ::get:: %s done", handler.Route)
		case	"POST":
			(*api).Post(handler.Route, _h)
			log.Printf("Request ::post:: %s done", handler.Route)
		case	"PUT":
			(*api).Put(handler.Route, _h)
			log.Printf("Request ::put:: %s done", handler.Route)
		case	"DELETE":
			(*api).Delete(handler.Route, _h)
			log.Printf("Request ::delete:: %s done", handler.Route)
		default:
			(*api).Get(handler.Route, _h)
			log.Printf("Request :: %s done", handler.Route)
		}
	}

}

func (repo *Repository) SetupRotes(app *fiber.App) {
	app.Use(logger.New())
	api := app.Group("/api")

	GenerateHandlers(&handlers.APIHandlers, &api, repo.DB)
}