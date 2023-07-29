package repository

import (
	"app/pkg/models"
	"context"
	"fmt"
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
		return fmt.Errorf("failed to Insert reservation menu: %w", err)
	}
	return nil
}
