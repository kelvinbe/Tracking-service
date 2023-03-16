package aft

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"tracking-service/utils"

	"github.com/samber/lo"
)

//Africa's Talking API logic

type AftClient struct {
	Username  string
	ApiKey    string
	ShortCode string
}

func NewAftClient() (*AftClient, error) {
	api_key := os.Getenv("AFRICAS_TALKING_API_KEY")
	username := os.Getenv("AFRICAS_TALKING_USERNAME")
	short_code := os.Getenv("AFRICAS_TALKING_SHORT_CODE")

	if lo.IsEmpty(username) || lo.IsEmpty(api_key) || lo.IsEmpty(short_code) {
		return nil, errors.New("AFT credentials not set")
	}

	return &AftClient{
		Username:  username,
		ApiKey:    api_key,
		ShortCode: short_code,
	}, nil
}

func (aft *AftClient) ActivateDevice(device_id string) error {
	client := &http.Client{}

	data := url.Values{}
	data.Set("from", aft.ShortCode)
	data.Set("to", "+254794699065")
	data.Set("message", utils.ACTIVATION_CODE)
	data.Set("username", aft.Username)

	encodedData := data.Encode()


	req, err := http.NewRequest(http.MethodPost, utils.AFRICAS_TALKING_SEND_URL, bytes.NewBufferString(encodedData))

	if err != nil {
		return err
	}

	req.Header.Add("apiKey", aft.ApiKey)
	req.Header.Add("username", aft.Username)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res , err := client.Do(req)

	if res.StatusCode != 200 && res.StatusCode != 201 && res.StatusCode != 202   {
		log.Printf("Error sending message: %v", res.StatusCode)
		error_str := fmt.Sprintf("error sending message %d", res.StatusCode)
		return errors.New(error_str)
	}

	if err != nil {
		log.Printf("Error sending message: %v", err)
		return err 
	}

	respBody, err := io.ReadAll(res.Body)


	if err != nil {
		log.Printf("Error reading response: %v", err)
		return err
	}

	log.Printf("Here is the response: %v", string(respBody))

	return nil
}

func (aft *AftClient) DeactivateDevice(device_id string) error {

	client := &http.Client{}

	data := url.Values{}

	data.Set("from", aft.ShortCode)
	data.Set("to", "+254794699065")
	data.Set("message", utils.DEACTIVATION_CODE)
	data.Set("username", aft.Username)

	encodedData := data.Encode()

	req, err := http.NewRequest(http.MethodPost, utils.AFRICAS_TALKING_SEND_URL, bytes.NewBufferString(encodedData))

	if err != nil {
		return err
	}

	req.Header.Add("apiKey", aft.ApiKey)
	req.Header.Add("username", aft.Username)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)

	if res.StatusCode != 200 && res.StatusCode != 201 && res.StatusCode != 202   {
		log.Printf("Error sending message: %v", res.StatusCode)
		error_str := fmt.Sprintf("error sending message %d", res.StatusCode)
		return errors.New(error_str)
	}

	if err != nil {
		log.Printf("Error sending message: %v", err)
		return err 
	}

	respBody, err := io.ReadAll(res.Body)

	if err != nil {
		log.Printf("Error reading response: %v", err)
		return err
	}

	log.Printf("Here is the response: %v", string(respBody))

	return nil
}
