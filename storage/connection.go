package storage

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/getsentry/sentry-go"
	lo "github.com/samber/lo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgressConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	SSLMode  string
}

type Connection struct {
	Client  *mongo.Client
	Context *context.Context
}

func NewMonongoClient() (*mongo.Client, error) {

	connection_string := os.Getenv("MONGO_CONNECTION_STRING")

	if lo.IsEmpty(connection_string) {
		return nil, errors.New("MONGO_CONNECTION_STRING is not set")
	}
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts :=  options.Client().ApplyURI(connection_string).SetServerAPIOptions(serverAPI)

	client, err := mongo.NewClient(opts)

	if err != nil {
		sentry.CaptureException(err)
		return nil, err
	}


	if err != nil {
		sentry.CaptureException(err)
		return nil, err
	}

	return client, nil

}

// this is for instantiating a postgres client
func NewPostgressClient(config PostgressConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", config.Host, config.Username, config.Password, config.Database, config.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		sentry.CaptureException(err)
		return nil, err
	}

	return db, nil
}
