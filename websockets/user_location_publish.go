package appwebsockets

import (
	"context"
	"log"
	"time"
	"tracking-service/dto"

	"github.com/getsentry/sentry-go"
	"github.com/gofiber/websocket/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func UserLocationPublish(app_clients *APP_CLIENTS, ws *websocket.Conn) interface{} {
	defer func() {
		if err := ws.Close(); err != nil {
			sentry.CaptureException(err)
			log.Printf("Failed to close websocket connection: %v", err)
		}
	}()

	for {
		log.Printf("User location publish websocket connected")
		user_location_info := &dto.IncomingUserLocationInfo{}

		err := ws.ReadJSON(user_location_info); if err != nil {
			sentry.CaptureException(err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Websocket error: %v", err)
			}
			return nil
		}

		log.Printf("User location info: %+v", user_location_info)

		_, err = app_clients.Mongo.Collection("user_location").InsertOne(context.TODO(), bson.M{
			"reservation_id": user_location_info.ReservationId,
			"latitude":       user_location_info.Latitude,
			"longitude":      user_location_info.Longitude,
			"timestamp":      time.Now(),
			"vehicle_id":    user_location_info.VehicleId,
		})

		if err != nil {
			sentry.CaptureException(err)
			log.Printf("An error occured while inserting user location: %v", err)
			ws.WriteMessage(websocket.TextMessage, []byte("ERROR"))
			continue
		}

		log.Printf("User location published")
		ws.WriteMessage(websocket.TextMessage, []byte("PUBLISHED"))
	}
	 // this is for simplicity on the client side, instead of sending a json response

	

}