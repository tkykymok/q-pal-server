package exmodels

import "github.com/volatiletech/null/v8"

type StaffWithRelated struct {
	StoreID            int         `boil:"store_id" json:"store_id" toml:"store_id" yaml:"store_id"`
	StaffID            int         `boil:"staff_id" json:"staff_id" toml:"staff_id" yaml:"staff_id"`
	Name               null.String `boil:"name" json:"name,omitempty" toml:"name" yaml:"name,omitempty"`
	Order              null.Int    `boil:"order" json:"order" toml:"order" yaml:"order"`
	BreakStartDatetime null.Time   `boil:"break_start_datetime" json:"break_start_datetime,omitempty" toml:"break_start_datetime" yaml:"break_start_datetime,omitempty"`
	BreakEndDatetime   null.Time   `boil:"break_end_datetime" json:"break_end_datetime,omitempty" toml:"break_end_datetime" yaml:"break_end_datetime,omitempty"`
	ActiveFlag         bool        `boil:"active_flag" json:"active_flag" toml:"active_flag" yaml:"active_flag"`
}
