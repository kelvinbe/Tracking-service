package storage

import (
	"context"
	"errors"
	"os"
	"time"

	lo "github.com/samber/lo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Connection struct {
	Client *mongo.Client
	Context *context.Context
}


func NewClient () (*Connection, error) {

	connection_string := os.Getenv("MONGO_CONNECTION_STRING")

	if lo.IsEmpty(connection_string) {
		return nil, errors.New("MONGO_CONNECTION_STRING is not set")
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(connection_string))

    if err != nil {
        return nil, err
    }

	ctx, _ := context.WithTimeout(context.Background(), 100 * time.Second)

	if err != nil {
		return nil, err
	}


	return &Connection{
		Client: client,
		Context: &ctx,
	} , nil;
	


}