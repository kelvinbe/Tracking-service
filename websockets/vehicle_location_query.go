package appwebsockets

import (
	"context"
	"log"
	"time"

	"tracking-service/dto"

	"github.com/getsentry/sentry-go"
	"github.com/gofiber/websocket/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func VehicleLocationQuery(app_clients *APP_CLIENTS, ws *websocket.Conn) interface{} {
	log.Printf("Vehicle location query websocket connected")
	defer func() {
		if err := ws.Close(); err != nil {
			sentry.CaptureException(err)
			log.Printf("Failed to close websocket connection: %v", err)
		}
	}()

	

	for {

		interval := time.Second * 5
		ticker := time.NewTicker(interval)

		vehicle_location_info := &dto.IncomingVehicleLocationInfo{}
	
		err := ws.ReadJSON(vehicle_location_info); if err != nil {
			sentry.CaptureException(err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Websocket error: %v", err)
			}
			return nil
		}

		go func(app_clients *APP_CLIENTS, ws *websocket.Conn) interface{} {
			for range ticker.C {

				cursor, err := app_clients.Mongo.Collection("user_location").Find(context.TODO(), bson.M{
					"vehicle_id": vehicle_location_info.VehicleId,
					"timestamp": bson.M{
						"$gte": time.Now().Add(-time.Minute * 5),
					},
				})

				if err != nil {
					sentry.CaptureException(err)
					if (err == mongo.ErrNoDocuments) {
						return ws.WriteJSON(bson.M{
							"message": "No locations found for the ID provided",
							"status":  "error",
							"data": nil,
						})
					}

					ws.WriteJSON(bson.M{
						"message": "No locations found",
						"status": "error",
						"data": nil,
					})
					break
				}

				var locations []bson.M


				err = cursor.All(context.TODO(), &locations); if err != nil {
					sentry.CaptureException(err)
					ws.WriteJSON(bson.M{
						"message": "No locations found",
						"status": "error",
						"data": nil,
					})
					break
				}

				ws.WriteJSON(bson.M{
					"data": locations,
					"message": "Locations found",
					"status": "success",
				})
			}
			return nil
		}(app_clients, ws)
	}

}