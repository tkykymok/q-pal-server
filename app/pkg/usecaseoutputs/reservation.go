package usecaseoutputs

import (
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/types"
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
	MenuID               int
	MenuName             null.String
	Price                types.NullDecimal
	Content              string
}

type WaitTime struct {
	ReservationNumber int
	Position          int
	Time              int
}

type ReservationIdentifyKey struct {
	ReservationID    int
	StoreID          int
	ReservedDatetime time.Time
}

type CreateReservation struct {
	ReservationNumber int
	Content           string
}
