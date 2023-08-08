package presenter

import (
	"app/pkg/enum"
	"app/pkg/usecaseoutputs"
	"app/pkg/utils"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/types"
	"net/url"
)

type Reservation struct {
	ReservationID        int               `json:"reservationId" `
	CustomerID           int               `json:"customerId" `
	StoreID              int               `json:"storeId" `
	StaffID              null.Int          `json:"staffId" `
	ReservationNumber    int               `json:"reservationNumber" `
	ReservedDatetime     string            `json:"reservedDatetime" `
	HoldStartDatetime    string            `json:"holdStartDatetime" `
	ServiceStartDatetime string            `json:"serviceStartDatetime" `
	ServiceEndDatetime   string            `json:"serviceEndDatetime" `
	Status               string            `json:"status" `
	ArrivalFlag          bool              `json:"arrivalFlag" `
	CancelType           null.Int          `json:"cancelType" `
	MenuID               int               `json:"menuId" `
	MenuName             string            `json:"menuName" `
	Price                types.NullDecimal `json:"price" `
	Content              string            `json:"content" `
}

type WaitTime struct {
	ReservationNumber int `json:"reservationNumber" `
	Position          int `json:"position"`
	Time              int `json:"time" `
}

type ReservationMessage struct {
	Message string `json:"message" `
}

type CreateReservation struct {
	ReservationNumber int    `json:"reservationNumber"`
	Content           string `json:"content"`
}

func GetReservationsResponse(data *[]usecaseoutputs.Reservation) ApiResponse {
	reservations := make([]Reservation, 0)
	for _, t := range *data {
		reservation := Reservation{
			ReservationID:        t.ReservationID,
			CustomerID:           t.CustomerID,
			StoreID:              t.StoreID,
			StaffID:              t.StaffID,
			ReservationNumber:    t.ReservationNumber,
			ReservedDatetime:     utils.ConvertTimeToString(t.ReservedDatetime),
			HoldStartDatetime:    utils.ConvertNTimeToString(t.HoldStartDatetime),
			ServiceStartDatetime: utils.ConvertNTimeToString(t.ServiceStartDatetime),
			ServiceEndDatetime:   utils.ConvertNTimeToString(t.ServiceEndDatetime),
			Status:               enum.ReservationStatusNames[enum.ReservationStatus(t.Status)],
			ArrivalFlag:          t.ArrivalFlag,
			CancelType:           t.CancelType,
			MenuID:               t.MenuID,
			MenuName:             utils.CheckString(t.MenuName),
			Price:                t.Price,
			Content:              url.QueryEscape(t.Content),
		}
		reservations = append(reservations, reservation)
	}

	return ApiResponse{
		Data:     reservations,
		Messages: []string{},
	}
}

func GetWaitTimeResponse(data *usecaseoutputs.WaitTime) ApiResponse {
	waitTime := WaitTime{
		ReservationNumber: data.ReservationNumber,
		Position:          data.Position,
		Time:              data.Time,
	}

	return ApiResponse{
		Data:     waitTime,
		Messages: []string{},
	}
}

func GetCreateReservationResponse(data *usecaseoutputs.CreateReservation, messages ...string) ApiResponse {
	encryptedStr := CreateReservation{
		ReservationNumber: data.ReservationNumber,
		Content:           url.QueryEscape(data.Content),
	}

	return ApiResponse{
		Data:     encryptedStr,
		Messages: messages,
	}
}
