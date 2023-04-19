package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	// "time"
	"tracking-service/repository"
	"tracking-service/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
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
	//setup 
	utils.LoadEnv()
	repo := repository.InitRepository() 
	app := fiber.New()
	repo.SetupRotes(app)

	// tests
	tests := []struct {
		description	string
		route string 
		expectedCode int
	}{
		{
			description: "Should return 200",
			route: "/api/polling?tracking_device_id=4438833",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest(http.MethodGet, test.route, nil)
		
		res, _ := app.Test(req)


		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)
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