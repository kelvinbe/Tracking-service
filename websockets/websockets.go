package appwebsockets

import (
	"context"
	"tracking-service/aft"

	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type APP_CLIENTS struct {
	Mongo        *mongo.Database
	Postgres     *gorm.DB
	MongoContext *context.Context
	Aft          *aft.AftClient
}