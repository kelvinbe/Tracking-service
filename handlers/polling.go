package handlers

import (
	"context"
	"net/http"
	"time"

	"tracking-service/dto"

	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)



func PollingHandler (app_client *APP_CLIENTS, ctx *fiber.Ctx) error {
	tracking_device_id := ctx.Query("tracking_device_id")


	if tracking_device_id == "" {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Reservation id is required",
			"status":  "error",
		})
	}

	var locations dto.TrackingDevice

	err := app_client.Mongo.Collection("tracking").FindOne(context.TODO(), bson.M{
		"tracking_device_id": tracking_device_id,
		"reservations": bson.M{
			"$elemMatch": bson.M{
				"status": bson.M{
					"$eq": "ACTIVE",
				},
			},
		},
		"reservations.locations": bson.M{
			"$elemMatch": bson.M{
				"time": bson.M{
					"$gt": time.Now().UTC().Add(-time.Minute * 5),
				},
			},
		},
	}, options.FindOne().SetProjection(bson.M{
		"reservations.$": 1,
		
		
	})).Decode(&locations)
	if err != nil {

		if (err == mongo.ErrNoDocuments) {
			return ctx.Status(http.StatusNotFound).JSON(&fiber.Map{
				"message": "No locations found",
				"status":  "error",
			})
		}
		sentry.CaptureException(err)
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Error while fetching locations",
			"status":  "error",
			"data":   err,
		})
	}
	
	if len(locations.Reservations) == 0 {
		return ctx.Status(http.StatusNotFound).JSON(&fiber.Map{
			"message": "No locations found",
			"status":  "error",
		})
	} 

	if len((locations.Reservations)[0].Locations) == 0 {
		return ctx.Status(http.StatusNotFound).JSON(&fiber.Map{
			"message": "No locations found",
			"status":  "error",
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Locations fetched successfully",
		"status":  "success",
		"data":    (locations.Reservations)[0].Locations,
	})
}