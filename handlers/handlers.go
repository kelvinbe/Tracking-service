package handlers

import (
	"context"
	"tracking-service/aft"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type APP_CLIENTS struct {
	Mongo        *mongo.Database
	Postgres     *gorm.DB
	MongoContext *context.Context
	Aft          *aft.AftClient
}

type APIHandler struct {
	Route   string
	Handler func(db *APP_CLIENTS, context *fiber.Ctx) error
	Method  string
}

func Ping(db *APP_CLIENTS, ctx *fiber.Ctx) error {
	return ctx.Status(200).JSON(&fiber.Map{
		"messsaage": "pong",
		"status":    "success",
	})
}

var APIHandlers []APIHandler = []APIHandler{
	{
		Route:   "/ping",
		Handler: Ping,
		Method:  "GET",
	},
}
