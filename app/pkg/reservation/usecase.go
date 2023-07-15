package reservation

import (
	"app/pkg/outputs"
	"context"
)

type Usecase interface {
	FetchReservationsByStoreId(ctx context.Context, storeId int) (*[]outputs.Reservation, error)
}

type usecase struct {
	repository Repository
}

func NewUsecase(r Repository) Usecase {
	return &usecase{
		repository: r,
	}
}

func (u usecase) FetchReservationsByStoreId(ctx context.Context, storeId int) (*[]outputs.Reservation, error) {
	reservations := make([]outputs.Reservation, 0)
	result, err := u.repository.ReadReservationsByStoreId(ctx, storeId)
	if err != nil {
		return nil, err
	}

	for _, t := range result {
		reservation := outputs.Reservation{
			ReservationID:        t.ReservationID,
			CustomerID:           t.CustomerID,
			StoreID:              t.StoreID,
			StaffID:              t.StaffID,
			ReservationNumber:    t.ReservationNumber,
			ReservedDatetime:     t.ReservedDatetime,
			HoldStartDatetime:    t.HoldStartDatetime,
			ServiceStartDatetime: t.ServiceStartDatetime,
			ServiceEndDatetime:   t.ServiceEndDatetime,
			Status:               t.Status,
			ArrivalFlag:          t.ArrivalFlag,
			CancelType:           t.CancelType,
		}
		reservations = append(reservations, reservation)
	}

	return &reservations, nil
}
