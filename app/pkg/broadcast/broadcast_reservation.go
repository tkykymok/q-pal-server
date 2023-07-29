package broadcast

import (
	"app/api/presenter"
	"app/pkg/enum"
)

var ReservationClient = NewReservationConnectionManager()

type ReservationConnectionManager struct {
	*ConnectionManager
}

func NewReservationConnectionManager() *ReservationConnectionManager {
	return &ReservationConnectionManager{
		ConnectionManager: NewConnectionManager(),
	}
}

func (manager *ReservationConnectionManager) SendNewReservation(subKey int, message presenter.ReservationMessage) {
	manager.SendMessage(enum.Reservation.String(), subKey, message)
}
