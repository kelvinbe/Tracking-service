package handlers

import (
	"context"
	"log"
	"net/http"
	"tracking-service/dto"

	"tracking-service/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
)



func LocationWebhookHandler (app_clients *APP_CLIENTS, ctx *fiber.Ctx) error {

	var requestBody dto.IncomingMessage;

	err := ctx.BodyParser(&requestBody)
	
	// send 200 to aft 
	ctx.SendStatus(http.StatusOK)

	if err != nil {
		// TODO: add a logger like sentry
		log.Printf("Error while parsing incoming message: %v", err)
		return err
	}

	lat, lon := utils.ExtractCoordsFromText(requestBody.Text)

	if lo.IsEmpty(lat) || lo.IsEmpty(lon) {
		log.Printf("Error while parsing incoming message: %v", err)
		return err
	}

	err = app_clients.Mongo.Collection("tracking").FindOneAndUpdate(context.TODO(), bson.M{
		"tracking_device_id": requestBody.From, // treating the number as the tracking device id
		"reservations.status": "ACTIVE",
	}, bson.M{
		"$push": bson.M{
			"reservations.locations": bson.M{
				"lat": lat,
				"lon": lon,
			},
		},
	}).Err()

	if err != nil {
		log.Printf("Error while updating location: %v", err)
		return err
	}


	return nil

}