package tests

import (
	"context"
	"testing"
	"tracking-service/repository"
	"tracking-service/utils"

	"github.com/gavv/httpexpect/v2"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func TestUserPublishLocation(t *testing.T) {
	utils.LoadEnv()
	// setup 
	reservation_id := uuid.NewV4().String()
	vehicle_id := uuid.NewV4().String()

	expectServer := httpexpect.Default(t, "http://0.0.0.0:8080")

	conn := expectServer.GET("/userlocation").WithWebsocketUpgrade().Expect()

	conn.Websocket().WriteJSON(map[string]interface{}{
		"reservation_id": reservation_id,
		"latitude":       "123.123",
		"longitude":      "123.123",
		"vehicle_id":     vehicle_id,
	})

	conn.Websocket().Expect().TextMessage().Body().IsEqual("PUBLISHED")

	conn.Websocket().Close()

	// teardown
	repo := repository.InitRepository()
	defer repo.Client.Disconnect(context.Background())
	repo.Mongo.Collection("user_location").DeleteMany(context.TODO(), bson.M{"reservation_id": reservation_id})

}

func TestVehicleLocationQuery(t *testing.T) {
	utils.LoadEnv()
	// setup
	reservation_id := uuid.NewV4().String()
	vehicle_id := uuid.NewV4().String()

	repo := repository.InitRepository()
	defer repo.Client.Disconnect(context.Background())

	_, err := repo.Mongo.Collection("user_location").InsertOne(context.TODO(), bson.M{
		"reservation_id": reservation_id,
		"latitude":       "123.123",
		"longitude":      "123.123",
		"vehicle_id":     vehicle_id,
	})

	if err != nil {
		t.Errorf("An error occured while inserting user location: %v", err)
	}


	expectServer := httpexpect.Default(t, "http://0.0.0.0:8080")

	conn := expectServer.GET("/vehiclelocation").WithWebsocketUpgrade().Expect()

	conn.Websocket().WriteJSON(map[string]interface{}{
		"vehicle_id": vehicle_id,
	})

	conn.Websocket().Expect().JSON().Object().ContainsKey("data").ContainsKey("message").ContainsKey("status")

	conn.Websocket().Close()

	// teardown
	repo.Mongo.Collection("user_location").DeleteMany(context.TODO(), bson.M{"reservation_id": reservation_id})

}

