package dto

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// these are data transfer objects, this file is mostly for defining the data structures that are used in the handlers

type FetchedReservation struct {
	ReservationId    uuid.UUID `json:"reservation_id"`
	VehicleId        uuid.UUID `json:"vehicle_id"`
	TrackingDeviceId string    `json:"tracking_device_id"`
}

type Location struct {
	Latitude  string    `bson:"lat"`
	Longitude string    `bson:"lon"`
	TimeStamp time.Time `bson:"time"` // not sure about the timestamp format, will make necessary modifications once sure
}

type MongoReservation struct {
	ReservationId string     `bson:"reservation_id"`
	Locations     []Location `bson:"locations"`
	Status        string     `bson:"status"`
}

type TrackingDevice struct {
	TrackingDeviceId string             `bson:"tracking_device_id"` // treating this as the device sim number
	Reservations     []MongoReservation `bson:"reservations"`
	Status           string             `bson:"status"`
}

// Africa's talking DTOs
type IncomingMessage struct {
	Date        string `json:"date"`
	From        string `json:"from"` // Treating this as the device sim number
	Id          string `json:"id"`
	LinkId      string `json:"linkId"`
	Text        string `json:"text"` // this is the message
	To          string `json:"to"`
	NetworkCode string `json:"networkCode"`
}

type Recipient struct {
	StatusCode int    `json:"statusCode"`
	Number     string `json:"number"`
	Cost       string `json:"cost"`
	MessageId  string `json:"messageId"`
}

type SMSMessageData struct {
	Message    *string     `json:"message"`
	Recepients []Recipient `json:"Recipients"`
}

type OutgoingMessageResp struct {
	SMSMessageData SMSMessageData `json:"SMSMessageData"`
}

type IncomingUserLocationInfo struct {
	VehicleId     string `json:"vehicle_id"`
	Longitude     string `json:"longitude"`
	Latitude      string `json:"latitude"`
	ReservationId string `json:"reservation_id"`
}

type IncomingVehicleLocationInfo struct {
	VehicleId string `json:"vehicle_id"`
}

type UserLocation struct {
	ReservationId string    `bson:"reservation_id" json:"reservation_id"`
	Latitude      string    `bson:"latitude" json:"latitude"`
	Longitude     string    `bson:"longitude" json:"longitude"`
	Timestamp     time.Time `bson:"timestamp" json:"timestamp"`
	VehicleId     string    `bson:"vehicle_id" json:"vehicle_id"` 
}



