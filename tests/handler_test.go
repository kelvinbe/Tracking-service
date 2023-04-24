package tests

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"tracking-service/repository"
	"tracking-service/utils"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestActivateHandler(t *testing.T) {

	// tests
	tests := []struct {
		description	string
		route string
		expectedCode int
	} {
		{

			description: "Should return 200",
			route: "/api/activate",
			expectedCode: http.StatusOK,
		},
	}

	// setup
	seed_data, err := Setup(); if err != nil {
		t.Error(err)
	}

	// run tests
	for _, test := range tests {
		url := fmt.Sprintf("http://0.0.0.0:8080%s", test.route)
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost,url , nil)
		if err != nil {
			t.Error(err)
		}
		res, err := client.Do(req)
		if err != nil {
			t.Error(err)
		}

		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)
	}

	// teardown
	err = Teardown(seed_data); if err != nil {
		t.Error(err)
	}
	t.Logf("Test completed successfully")
}

func TestPollingHandler(t *testing.T){
	utils.LoadEnv()
	with_no_location := uuid.NewV4().String()

	with_four_locations := uuid.NewV4().String()


	// setup
	repo := repository.InitRepository()

	// add a tracking device with no location
	_, err := repo.Mongo.Collection("tracking_devices").InsertOne(context.TODO(), bson.M{
		"tracking_device_id": with_no_location,
		"status": "ACTIVE",
		"reservations": []bson.M{
			{
				"reservation_id": uuid.NewV4().String(),
				"locations": []bson.M{},
				"status": "ACTIVE",
			},
		},
	}); if err != nil {
		t.Errorf("An error occured while inserting tracking device: %v", err)
	}

	// and one with four locations
	_, err = repo.Mongo.Collection("tracking_devices").InsertOne(context.TODO(), bson.M{
		"tracking_device_id": with_four_locations,
		"status": "ACTIVE",
		"reservations": []bson.M{
			{
				"reservation_id": uuid.NewV4().String(),
				"locations": []bson.M{
				},
				"status": "ACTIVE",
			},
		},
	}); if err != nil {
		t.Errorf("An error occured while inserting tracking device: %v", err)
	}



	if err != nil {
		t.Errorf("An error occured while inserting tracking device: %v", err)
	}

	// tests
	tests := []struct {
		description	string
		route string 
		expectedCode int
	}{
		{
			description: "Should return 404",
			route: fmt.Sprintf("/api/polling?tracking_device_id=%s", with_no_location),
			expectedCode: http.StatusNotFound,
		},
		{
			description: "Should return 200",
			route: fmt.Sprintf("/api/polling?tracking_device_id=%s", with_four_locations),
			expectedCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		
		url := fmt.Sprintf("http://0.0.0.0:8080%s", test.route)
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet,url , nil)
		if err != nil {
			t.Error(err)
		}

		res, err := client.Do(req)

		if err != nil {
			t.Error(err)
		}

		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)
		
	}

	// teardown
	for _, tracking_device_id := range []string{with_no_location, with_four_locations} {
		_, err := repo.Mongo.Collection("tracking_devices").DeleteOne(context.TODO(), bson.M{
			"tracking_device_id": tracking_device_id,
		}); if err != nil {
			t.Errorf("An error occured while deleting tracking device: %v", err)
		}
	}

}

func TestDeactivateHandler(t *testing.T) {
	//setup 
	seed_data, err := Setup(); if err != nil {
		t.Error(err)
	}

	// activate tracking devices
	url := fmt.Sprintf("http://0.0.0.0:8080%s", "/api/activate")
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost,url , nil)
	if err != nil {
		t.Error(err)
	}
	_, err = client.Do(req)
	if err != nil {
		t.Error(err)
	}

	// tests { these have to be added in chronological order i.e deactivate, whatever happens after that, etc}
	tests := []struct {
		description	string
		route string 
		expectedCode int
	}{
		{
			description: "Should return 200",
			route: "/api/deactivate",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		test_url := fmt.Sprintf("http://0.0.0.0:8080%s", test.route)
		req, err := http.NewRequest(http.MethodPost,test_url , nil)

		if err != nil {
			t.Error(err)
		}
		res, err := client.Do(req)

		if err != nil {
			t.Error(err)
		}

		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)
	}

	// teardown
	err = Teardown(seed_data); if err != nil {
		t.Error(err)
	}
	t.Logf("Test completed successfully")
}