package tests

import (
	"context"
	"fmt"
	"time"
	"tracking-service/repository"
	"tracking-service/utils"

	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type SeedDataIds struct {
	Market string 
	SubMarket string
	Hosts []string
	Stations []string
	Customers []string 
	Reservations []string
	Repo *repository.Repository
}


func Setup () (*SeedDataIds, error) {
	//setup 
	utils.LoadEnv()

	repo := repository.InitRepository() 
	// market ids
	ken_id := uuid.NewV4().String()

	// submarket ids
	ken_nai_id := uuid.NewV4().String()

	// create markets
	err := repo.Postgres.Exec(fmt.Sprintf(`
		insert into "public"."Market" (id,country, name, currency)
		values ('%s','Kenya', 'Nairobi', 'KES');
	`, ken_id)).Error; if err != nil {
		return nil, err
	}

	

	// create submarkets
	err = repo.Postgres.Exec(fmt.Sprintf(`
		insert into "public"."SubMarket" (id,market_id, name)
		values ('%s','%s', 'Nairobi');
	`, ken_nai_id, ken_id)).Error; if err != nil {
		return nil, err
	}

	var host_ids = []string{
		uuid.NewV4().String(),
		uuid.NewV4().String(),
		uuid.NewV4().String(),
		uuid.NewV4().String(),
	}

	// create host users
	for _, host_id := range host_ids {
		host_error := repo.Postgres.Exec(fmt.Sprintf(`
			insert into "public"."User" (id, sub_market_id, email, handle, uid, user_type)
			values ('%s', '%s', '%s', '%s', '%s', 'HOST');
		`, host_id, ken_nai_id, fmt.Sprintf("%s@email.com", host_id), host_id, host_id)).Error; if host_error != nil {
			return nil, host_error
		}
	}

	// create stations
	var station_ids = []string{
		uuid.NewV4().String(),
		uuid.NewV4().String(),
		uuid.NewV4().String(),
		uuid.NewV4().String(),
	}

	for i, station_id := range station_ids {
		station_err := repo.Postgres.Exec(`
			insert into "public"."Station" (id, name, sub_market_id, user_id)
			values (?, ?, ?, ?);
		`, station_id, station_id, ken_nai_id, host_ids[i]).Error; if station_err != nil {
			return nil, station_err
		}
	}

	// create vehicles with tracking device ids

	var vehicle_ids = []string{
		uuid.NewV4().String(),
		uuid.NewV4().String(),
		uuid.NewV4().String(),
		uuid.NewV4().String(),
	}

	var tracking_device_ids = []string{ // these are just phone numbers for those devices
		"1234567",
		"1234568",
		"1234569",
		"1234570",
	}

	for i, vehicle_id := range vehicle_ids {
		repo.Postgres.Exec(`
			insert into "public"."Vehicle" (id, user_id, station_id, tracking_device_id)
			values (?, ?, ?, ?);
		`, vehicle_id, host_ids[i], station_ids[i], tracking_device_ids[i])
	}

	
	// create customers for reservations
	var customer_ids = []string{
		uuid.NewV4().String(),
		uuid.NewV4().String(),
		uuid.NewV4().String(),
		uuid.NewV4().String(),
	}

	for _, customer_id := range customer_ids {
		repo.Postgres.Exec(`
			insert into "public"."User" (id, sub_market_id, email, handle, uid, user_type)
			values (?, ?, ?, ?, ?, 'CUSTOMER');
		`, customer_id, ken_nai_id, fmt.Sprintf("%s@email.com", customer_id), customer_id, customer_id)
	}

	// create reservations for the customers
	var reservation_ids = []string{
		uuid.NewV4().String(),
		uuid.NewV4().String(),
		uuid.NewV4().String(),
		uuid.NewV4().String(),
	}

	for i, reservation_id := range reservation_ids {
		repo.Postgres.Exec(`
			insert into "public"."Reservation" (id, user_id, vehicle_id, start_date_time, end_date_time, status, updated_at)
			values (?, ?, ?, ?, ?, 'ACTIVE', ?);
		`, reservation_id, customer_ids[i], vehicle_ids[i], (time.Now().Add( - time.Hour * 5 )), time.Now().Add(time.Hour * 2), time.Now())
	}


	return &SeedDataIds{
		Market: ken_id,
		SubMarket: ken_nai_id,
		Hosts: host_ids,
		Stations: station_ids,
		Customers: customer_ids,
		Reservations: reservation_ids,
		Repo: repo,
	}, nil
}

func Teardown (data *SeedDataIds) error {

	// delete reservations
	for _, reservation_id := range data.Reservations {
		err := data.Repo.Postgres.Exec(`
			delete from "public"."Reservation" where id = ?
		`, reservation_id).Error; if err != nil {
			return err
		}
	}

	// delete customers
	for _, customer_id := range data.Customers {
		err := data.Repo.Postgres.Exec(`
			delete from "public"."User" where id = ?
		`, customer_id).Error; if err != nil {
			return err
		}
	}

	// delete vehicles
	for _, vehicle_id := range data.Reservations {
		err := data.Repo.Postgres.Exec(`
			delete from "public"."Vehicle" where id = ?
		`, vehicle_id).Error; if err != nil {
			return err;
		}
	}

	// delete stations
	for _, station_id := range data.Stations {
		vehicles_delete_err := data.Repo.Postgres.Exec(`
			delete from "public"."Vehicle" where station_id = ?
		`, station_id).Error; if vehicles_delete_err != nil {
			return vehicles_delete_err
		}
		err := data.Repo.Postgres.Exec(`
			delete from "public"."Station" where id = ?
		`, station_id).Error; if err != nil {
			return err
		}
	}

	// delete hosts
	for _, host_id := range data.Hosts {
		err := data.Repo.Postgres.Exec(`
			delete from "public"."User" where id = ?
		`, host_id).Error; if err != nil {
			return err
		}
	}

	// delete submarkets
	submarket_err := data.Repo.Postgres.Exec(`
		delete from "public"."SubMarket" where id = ?
	`, data.SubMarket).Error; if submarket_err != nil {
		return submarket_err
	}

	// delete markets
	market_err := data.Repo.Postgres.Exec(`
		delete from "public"."Market" where id = ?
	`, data.Market).Error; if market_err != nil {
		return market_err
	}

	// delete all tracking device documents from mongo
	_, err := data.Repo.Mongo.Collection("tracking").DeleteMany(context.Background(), bson.M{
		"tracking_device_id": bson.M{
			"$in": []string{
				"1234567",
				"1234568",
				"1234569",
				"1234570",
			},
		},
	})

	if err != nil {
		return err
	}

	cerr := data.Repo.Mongo.Client().Disconnect(context.Background()); if cerr != nil {
		return cerr
	}

	return nil
}



