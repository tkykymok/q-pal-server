package usecaseinputs

import "github.com/volatiletech/null/v8"

type CreateReservationInput struct {
	CustomerID int
	StoreID    int
	MenuID     int
}

type UpdateReservationInput struct {
	ReservationID int      `json:"reservationId"`
	StoreID       int      `json:"storeId"`
	Status        string   `json:"status"`
	StaffID       null.Int `json:"staffId"`
	MenuID        null.Int `json:"menuId"`
}
