package usecaseoutputs

import (
	"github.com/volatiletech/null/v8"
)

type Staff struct {
	StoreID            int
	StaffID            int
	Name               null.String
	Order              null.Int
	BreakStartDatetime null.Time
	BreakEndDatetime   null.Time
	ActiveFlag         bool
}
