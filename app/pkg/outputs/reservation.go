package outputs

import (
	"github.com/volatiletech/null/v8"
	"time"
)

type Reservation struct {
	ReservationID        int       `json:"reservation_id" `
	CustomerID           int       `json:"customer_id" `
	StoreID              int       `json:"store_id" `
	StaffID              null.Int  `json:"staff_id" `
	ReservationNumber    int       `json:"reservation_number" `
	ReservedDatetime     time.Time `json:"reserved_datetime" `
	HoldStartDatetime    null.Time `json:"hold_start_datetime" `
	ServiceStartDatetime null.Time `json:"service_start_datetime" `
	ServiceEndDatetime   null.Time `json:"service_end_datetime" `
	Status               null.Int  `json:"status" `
	ArrivalFlag          bool      `json:"arrival_flag" `
	CancelFlag           bool      `json:"cancel_flag" `
	CancelType           null.Int  `json:"cancel_type" `
}
