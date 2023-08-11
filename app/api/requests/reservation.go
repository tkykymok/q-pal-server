package requests

import "github.com/volatiletech/null/v8"

type GetReservations struct {
	StoreId int
}

type CreateReservation struct {
	StoreId int
}

type UpdateReservation struct {
	ReservationID int      `json:"reservationId"`
	Status        string   `json:"status"`
	StaffID       null.Int `json:"staffId"`
	MenuID        null.Int `json:"menuId"`
}
