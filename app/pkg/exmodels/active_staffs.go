package exmodels

import "github.com/volatiletech/null/v8"

type AvailableStaff struct {
	StaffID            int       `boil:"staff_id" json:"staff_id" toml:"staff_id" yaml:"staff_id"`
	StoreID            int       `boil:"store_id" json:"store_id,omitempty" toml:"store_id" yaml:"store_id,omitempty"`
	BreakStartDatetime null.Time `boil:"break_start_datetime" json:"break_start_datetime,omitempty" toml:"break_start_datetime" yaml:"break_start_datetime,omitempty"`
	BreakEndDatetime   null.Time `boil:"break_end_datetime" json:"break_end_datetime,omitempty" toml:"break_end_datetime" yaml:"break_end_datetime,omitempty"`
	ReservationID      null.Int  `boil:"reservation_id" json:"reservation_id" toml:"reservation_id" yaml:"reservation_id"`
}
