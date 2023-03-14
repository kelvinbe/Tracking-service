package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type APIHandler struct {
	Route		string
	Handler		func(db *mongo.Database, context *fiber.Ctx) error
	Method		string
}

func Ping (db *mongo.Database, ctx *fiber.Ctx) error {
	return ctx.Status(200).JSON(&fiber.Map{
		"messsaage": "pong",
		"status": "success",
	})
}

var APIHandlers []APIHandler = []APIHandler{
	{
		Route: "/ping",
		Handler: Ping,
		Method: "GET",
	},
}