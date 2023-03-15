package dto

import uuid "github.com/satori/go.uuid"

// these are data transfer objects, this file is mostly for defining the data structures that are used in the handlers

type FetchedReservation struct {
	ReservationId    uuid.UUID `json:"reservation_id"`
	VehicleId        uuid.UUID `json:"vehicle_id"`
	TrackingDeviceId string    `json:"tracking_device_id"`
}

type Location struct {
	Latitude  float64 `bson:"latitude"`
	Longitude float64 `bson:"longitude"`
	TimeStamp int64   `bson:"timestamp"` // not sure about the timestamp format, will make necessary modifications once sure
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
