package usecaseoutputs

import (
	"github.com/volatiletech/null/v8"
	"time"
)

type Reservation struct {
	ReservationID        int
	CustomerID           int
	StoreID              int
	StaffID              null.Int
	ReservationNumber    int
	ReservedDatetime     time.Time
	HoldStartDatetime    null.Time
	ServiceStartDatetime null.Time
	ServiceEndDatetime   null.Time
	Status               int
	ArrivalFlag          bool
	CancelType           null.Int
}

type WaitTime struct {
	ReservationNumber int
	Position          int
	Time              int
}
