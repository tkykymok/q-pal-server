package repository

import (
	"app/api/errpkg"
	"app/pkg/models"
	"context"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type ReservationMenuRepository interface {
	InsertReservationMenu(ctx context.Context, reservationMenu *models.ReservationMenu) error
}

type reservationMenuRepository struct {
}

func NewReservationMenuRepo() ReservationMenuRepository {
	return &reservationMenuRepository{}
}

func (r reservationMenuRepository) InsertReservationMenu(ctx context.Context, reservationMenu *models.ReservationMenu) error {
	err := reservationMenu.InsertG(ctx, boil.Infer())
	if err != nil {
		return &errpkg.DatabaseError{
			InternalError: err,
			Operation:     "InsertReservationMenu",
		}
	}
	return nil
}
