package exmodels

import (
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/types"
	"time"
)

type ReservationWithRelated struct {
	ReservationID        int               `boil:"reservation_id" json:"reservation_id" toml:"reservation_id" yaml:"reservation_id"`
	CustomerID           int               `boil:"customer_id" json:"customer_id" toml:"customer_id" yaml:"customer_id"`
	StoreID              int               `boil:"store_id" json:"store_id" toml:"store_id" yaml:"store_id"`
	StaffID              null.Int          `boil:"staff_id" json:"staff_id,omitempty" toml:"staff_id" yaml:"staff_id,omitempty"`
	ReservationNumber    int               `boil:"reservation_number" json:"reservation_number" toml:"reservation_number" yaml:"reservation_number"`
	ReservedDatetime     time.Time         `boil:"reserved_datetime" json:"reserved_datetime" toml:"reserved_datetime" yaml:"reserved_datetime"`
	HoldStartDatetime    null.Time         `boil:"hold_start_datetime" json:"hold_start_datetime,omitempty" toml:"hold_start_datetime" yaml:"hold_start_datetime,omitempty"`
	ServiceStartDatetime null.Time         `boil:"service_start_datetime" json:"service_start_datetime,omitempty" toml:"service_start_datetime" yaml:"service_start_datetime,omitempty"`
	ServiceEndDatetime   null.Time         `boil:"service_end_datetime" json:"service_end_datetime,omitempty" toml:"service_end_datetime" yaml:"service_end_datetime,omitempty"`
	Status               int               `boil:"status" json:"status,omitempty" toml:"status" yaml:"status,omitempty"`
	ArrivalFlag          bool              `boil:"arrival_flag" json:"arrival_flag" toml:"arrival_flag" yaml:"arrival_flag"`
	CancelType           null.Int          `boil:"cancel_type" json:"cancel_type,omitempty" toml:"cancel_type" yaml:"cancel_type,omitempty"`
	MenuID               int               `boil:"menu_id" json:"menu_id" toml:"menu_id" yaml:"menu_id"`
	MenuName             null.String       `boil:"menu_name" json:"menu_name,omitempty" toml:"menu_name" yaml:"menu_name,omitempty"`
	Price                types.NullDecimal `boil:"price" json:"price,omitempty" toml:"price" yaml:"price,omitempty"`
	Time                 int               `boil:"time" json:"time,omitempty" toml:"time" yaml:"time,omitempty"`
}

type HandleTime struct {
	CustomerID int `boil:"customer_id" json:"customer_id" toml:"customer_id" yaml:"customer_id"`
	StoreID    int `boil:"store_id" json:"store_id" toml:"store_id" yaml:"store_id"`
	MenuID     int `boil:"menu_id" json:"menu_id" toml:"menu_id" yaml:"menu_id"`
	Time       int `boil:"time" json:"time,omitempty" toml:"time" yaml:"time,omitempty"`
}
