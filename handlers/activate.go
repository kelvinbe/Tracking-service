package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"tracking-service/dto"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func ActivateHandler(app_clients *APP_CLIENTS, ctx *fiber.Ctx) error {

	var activeReservations *[]dto.FetchedReservation
	var result *string

	err := app_clients.Postgres.Raw(`
		select json_agg(json_build_object(
			'reservation_id', res.id,
			'vehicle_id', veh.id,
			'tracking_device_id', veh.tracking_device_id
		)) as result from "public"."Reservation" as res 
		inner join "public"."Vehicle" as veh on res.vehicle_id = veh.id
		where res.status in('ACTIVE')
		and res.type != 'BLOCK'
		and veh.tracking_device_id is not null;
	`).Scan(&result).Error

	log.Printf("result: %v", *result)

	if result == nil {
		return ctx.Status(http.StatusNotFound).JSON(&fiber.Map{
			"message": "No reservation to activate",
			"status":  "error",
		})
	}

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Error while fetching active reservations",
			"status":  "error",
		})
	}

	err = json.Unmarshal([]byte(*result), &activeReservations)

	log.Printf("activeReservations: %v", *activeReservations)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Error while parsing active reservations",
			"status":  "error",
			"error": err,
		})
	}

	var loop_errors []error
	for _, activeReservation := range *activeReservations {
		log.Printf("activeReservation: %v", activeReservation)
		var reservation = activeReservation
		var trackingDevice *dto.TrackingDevice
		err := app_clients.Mongo.Collection("tracking").FindOne(context.TODO(), bson.M{
			"tracking_device_id": reservation.TrackingDeviceId,
		}).Decode(&trackingDevice)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				var newTrackingDevice = dto.TrackingDevice{
					TrackingDeviceId: reservation.TrackingDeviceId,
					Reservations:     []dto.MongoReservation{},
					Status:           "INACTIVE",
				}
				_, err = app_clients.Mongo.Collection("tracking").InsertOne(context.TODO(), newTrackingDevice)

				if err != nil {
					loop_errors = append(loop_errors, err)
					// continue to the next active reservation
					continue
				}

				// also add in the reservation
				var newReservation = dto.MongoReservation{
					ReservationId: reservation.ReservationId.String(),
					Locations:     []dto.Location{},
					Status:        "INACTIVE",
				}

				log.Printf("Here is the new reservation %v", newReservation)
				_, err = app_clients.Mongo.Collection("tracking").UpdateOne(context.TODO(), bson.M{
					"tracking_device_id": reservation.TrackingDeviceId,
				}, bson.M{
					"$push": bson.M{
						"reservations": newReservation,
					},
				})

				if err != nil {
					loop_errors = append(loop_errors, err)
					continue
				}
				log.Printf("No issue adding the new reservation")
			} else {
				loop_errors = append(loop_errors, err)
				continue
			}
		}

		// find the tracking device again after adding it
		err = app_clients.Mongo.Collection("tracking").FindOne(context.TODO(), bson.M{
			"tracking_device_id": reservation.TrackingDeviceId,
		}).Decode(&trackingDevice)

		if err != nil {
			loop_errors = append(loop_errors, err)
			continue
		}

		// if the tracking device is inactive, activate it
		if trackingDevice.Status == "INACTIVE" {
			// activate the tracking device
			err = app_clients.Aft.ActivateDevice(trackingDevice.TrackingDeviceId)

			if err != nil {
				loop_errors = append(loop_errors, err)
				continue
			}

			//device is now active, update the status in the database
			_, err = app_clients.Mongo.Collection("tracking").UpdateOne(context.TODO(), bson.M{
				"tracking_device_id": reservation.TrackingDeviceId,
			}, bson.M{
				"$set": bson.M{
					"status": "ACTIVE",
				},
			})

			if err != nil {
				loop_errors = append(loop_errors, err)
				continue
			}

			// also update the status of the reservation
			_, err = app_clients.Mongo.Collection("tracking").UpdateOne(context.TODO(), bson.M{
				"tracking_device_id":          reservation.TrackingDeviceId,
				"reservations.reservation_id": reservation.ReservationId.String(),
			}, bson.M{
				"$set": bson.M{
					"reservations.$.status": "ACTIVE",
				},
			})

			if err != nil {
				loop_errors = append(loop_errors, err)
				continue
			}
		}
		// else just continue to the next active reservation

		continue

	}

	if len(loop_errors) > 0 {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Error while activating tracking devices",
			"status":  "error",
			"errors":  loop_errors,
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Tracking devices activated successfully",
		"status":  "success",
	})
}
