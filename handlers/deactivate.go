package handlers

import (
	"encoding/json"
	"net/http"
	"tracking-service/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
)

func DeactivateDevice(app_clients *APP_CLIENTS, ctx *fiber.Ctx) error {

	var inactive_reservations []dto.FetchedReservation
	var results *string
	err := app_clients.Postgres.Raw(`
		select res.id as reservation_id, veh.id as vehicle_id, veh.tracking_device_id as tracking_device_id  from Reservation as res
		inner join Vehicle as veh on res.vehicle_id = veh.id
		where status in('CANCELLED', 'COMPLETE')
		and type != 'BLOCK'
	`).Scan(&results).Error

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Error while fetching inactive reservations",
			"status":  "error",
		})
	}

	err = json.Unmarshal([]byte(*results), &inactive_reservations)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Error while parsing inactive reservations",
			"status":  "error",
		})
	}

	var loop_errors []error
	for _, inactive_reservation := range inactive_reservations {
		var reservation = inactive_reservation

		// get the tracking device
		var trackingDevice dto.TrackingDevice
		err := app_clients.Mongo.Collection("tracking").FindOne(*app_clients.MongoContext, bson.M{
			"tracking_device_id": reservation.TrackingDeviceId,
		}).Decode(&trackingDevice)

		if err != nil {
			// unlike activate if we don't get the tracking device we just skip it
			loop_errors = append(loop_errors, err)
			continue
		}

		if trackingDevice.Status == "ACTIVE" {
			// if tracking device is active we can deactivate it

			err := app_clients.Aft.DeactivateDevice(trackingDevice.TrackingDeviceId)

			if err != nil {
				// will try again on the next request so we just skip it
				loop_errors = append(loop_errors, err)
				continue
			}

			// update the tracking device status
			_, err = app_clients.Mongo.Collection("tracking").UpdateOne(*app_clients.MongoContext, bson.M{
				"tracking_device_id": trackingDevice.TrackingDeviceId,
			}, bson.M{
				"$set": bson.M{
					"status": "INACTIVE",
				},
			})

			// update the one active reservation
			the_active_reservation, found := lo.Find(trackingDevice.Reservations, func (reservation dto.MongoReservation) bool {
				return reservation.Status == "ACTIVE"
			})

			if found {
				_, err = app_clients.Mongo.Collection("tracking").UpdateOne(*app_clients.MongoContext, bson.M{
					"tracking_device_id": trackingDevice.TrackingDeviceId,
					"reservations.reservation_id": the_active_reservation.ReservationId,
				}, bson.M{
					"$set": bson.M{
						"reservations.$.status": "INACTIVE",
					},
				})

				if err != nil {
					// will try again on the next request so we just skip it
					loop_errors = append(loop_errors, err)
					continue
				}
			}

			if err != nil {
				// will try again on the next request so we just skip it
				loop_errors = append(loop_errors, err)
				continue
			}



		}

	}

	if len(loop_errors) > 0 {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Error while deactivating tracking devices",
			"status":  "error",
			"errors":  loop_errors,
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Tracking devices deactivated successfully",
		"status":  "success",
	})

}
