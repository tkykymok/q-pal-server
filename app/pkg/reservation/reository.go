package reservation

import (
	"app/pkg/models"
	"context"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"time"
)

type Repository interface {
	ReadReservationsByStoreId(ctx context.Context, storeId int) (models.ReservationSlice, error)
}

type repository struct {
}

func NewRepo() Repository {
	return &repository{}
}

func (r repository) ReadReservationsByStoreId(ctx context.Context, storeId int) (models.ReservationSlice, error) {
	currentTime := time.Now()

	// 現在の日付の開始時刻を取得
	startOfDay := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
	// 現在の日付の終了時刻を取得
	endOfDay := startOfDay.Add(time.Hour*24 - time.Second)

	mods := []qm.QueryMod{
		qm.Where("store_id = ?", storeId),
		qm.Where("cancel_flag = ?", 0),
		qm.Where("reserved_datetime >= ?", startOfDay),
		qm.Where("reserved_datetime <= ?", endOfDay),
	}

	return models.Reservations(mods...).AllG(ctx)
}
