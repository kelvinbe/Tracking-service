package storage

import (
	"context"
	"errors"
	"fmt"
	"os"

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

	client, err := mongo.NewClient(options.Client().ApplyURI(connection_string))

	if err != nil {
		return nil, err
	}


	if err != nil {
		return nil, err
	}

	return client, nil

}

// this is for instantiating a postgres client
func NewPostgressClient(config PostgressConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", config.Host, config.Username, config.Password, config.Database, config.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil
}
