package repository

import (
	"context"
	"log"
	"os"
	"tracking-service/aft"
	"tracking-service/handlers"
	"tracking-service/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type Repository struct {
	Client   *mongo.Client
	Mongo    *mongo.Database
	Postgres *gorm.DB
	Aft      *aft.AftClient // Africa's Talking API client
}

func InitRepository() *Repository {

	client, err := storage.NewMonongoClient()

	if err != nil {
		log.Fatalf("Error connecting to mongo database: %s", err.Error())
	}

	err = client.Connect(context.Background())

	if err != nil {
		log.Fatalf("Error connecting to mongo database: %s", err.Error())
	}

	db := client.Database("tracking_service_test")

	postgres, err := storage.NewPostgressClient(storage.PostgressConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	})

	if err != nil {
		log.Fatalf("Error connecting to postgres database: %s", err.Error())
	}

	aft_client, err := aft.NewAftClient()

	if err != nil {
		log.Fatalf("Error initializing AFT client: %s", err.Error())
	}

	return &Repository{
		Client:   client,
		Mongo:    db,
		Postgres: postgres,
		Aft:      aft_client,
	}
}

// function to make it easier to add routes
func GenerateHandlers(handlers *[]handlers.APIHandler, api *fiber.Router, db *handlers.APP_CLIENTS) {

	for i := range *handlers {
		handler := (*handlers)[i]
		_h := func(ctx *fiber.Ctx) error {
			return handler.Handler(db, ctx)
		}
		switch handler.Method {
		case "GET":
			(*api).Get(handler.Route, _h)
			log.Printf("Request ::get:: %s done", handler.Route)
		case "POST":
			(*api).Post(handler.Route, _h)
			log.Printf("Request ::post:: %s done", handler.Route)
		case "PUT":
			(*api).Put(handler.Route, _h)
			log.Printf("Request ::put:: %s done", handler.Route)
		case "DELETE":
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

	GenerateHandlers(&handlers.APIHandlers, &api, &handlers.APP_CLIENTS{
		Mongo:        repo.Mongo,
		Postgres:     repo.Postgres,
		Aft:          repo.Aft,
	})
}
