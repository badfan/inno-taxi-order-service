package models

import (
	"time"

	"github.com/google/uuid"
)

type TaxiType string

const (
	TaxiTypeEconomy  TaxiType = "economy"
	TaxiTypeComfort  TaxiType = "comfort"
	TaxiTypeBusiness TaxiType = "business"
	TaxiTypeElectro  TaxiType = "electro"
)

type Order struct {
	ID          int32     `json:"id"`
	UserUuid    uuid.UUID `json:"user_uuid"`
	DriverUuid  uuid.UUID `json:"driver_uuid"`
	Origin      string    `json:"origin"`
	Destination string    `json:"destination"`
	TaxiType    TaxiType  `json:"taxi_type"`
	IsFinished  bool      `json:"is_finished"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
