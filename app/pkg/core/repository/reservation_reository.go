package repository

import (
	"app/pkg/enum"
	"app/pkg/exmodels"
	"app/pkg/models"
	"context"
	"fmt"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"strings"
	"time"
)

type ReservationRepository interface {
	ReadTodayReservations(ctx context.Context, storeId int, status ...enum.ReservationStatus) (*[]exmodels.ReservationWithRelated, error)
	ReadHandleTimes(ctx context.Context, storeId int) (*[]exmodels.HandleTime, error)
}

type reservationRepository struct {
}

func NewReservationRepo() ReservationRepository {
	return &reservationRepository{}
}

func (r reservationRepository) ReadTodayReservations(ctx context.Context, storeId int, status ...enum.ReservationStatus) (*[]exmodels.ReservationWithRelated, error) {
	// SELECTするカラム
	selectCols := []string{
		models.ReservationTableColumns.ReservationID,
		models.ReservationTableColumns.CustomerID,
		models.ReservationTableColumns.StoreID,
		models.ReservationTableColumns.StaffID,
		models.ReservationTableColumns.ReservationNumber,
		models.ReservationTableColumns.ReservedDatetime,
		models.ReservationTableColumns.HoldStartDatetime,
		models.ReservationTableColumns.ServiceStartDatetime,
		models.ReservationTableColumns.ServiceEndDatetime,
		models.ReservationTableColumns.Status,
		models.ReservationTableColumns.ArrivalFlag,
		models.ReservationTableColumns.CancelType,
		models.ReservationMenuTableColumns.MenuID,
		models.MenuTableColumns.MenuName,
		models.MenuTableColumns.Price,
		models.MenuTableColumns.Time,
	}

	currentTime := time.Now()
	// 現在の日付の開始時刻を取得
	startOfDay := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
	// 現在の日付の終了時刻を取得
	endOfDay := startOfDay.Add(time.Hour*24 - time.Second)

	// Convert ReservationStatus slice to int slice
	statusInts := make([]interface{}, len(status))
	for i, s := range status {
		statusInts[i] = int(s)
	}

	// QueryModの生成
	mods := []qm.QueryMod{
		qm.Select(strings.Join(selectCols, ",")),
		qm.InnerJoin(fmt.Sprintf("%s on %s = %s",
			models.TableNames.ReservationMenus,
			models.ReservationTableColumns.ReservationID,
			models.ReservationMenuTableColumns.ReservationID,
		)),
		qm.InnerJoin(fmt.Sprintf("%s on %s = %s and %s = %s",
			models.TableNames.Menus,
			models.ReservationMenuTableColumns.MenuID,
			models.MenuTableColumns.MenuID,
			models.ReservationMenuTableColumns.StoreID,
			models.MenuTableColumns.StoreID,
		)),
		qm.Where(fmt.Sprintf("%s = ?", models.ReservationTableColumns.StoreID), storeId),
		qm.Where(fmt.Sprintf("%s >= ?", models.ReservationTableColumns.ReservedDatetime), startOfDay),
		qm.Where(fmt.Sprintf("%s <= ?", models.ReservationTableColumns.ReservedDatetime), endOfDay),
		qm.WhereIn(fmt.Sprintf("%s in ?", models.ReservationTableColumns.Status), statusInts...),
		qm.OrderBy(fmt.Sprintf("%s", models.ReservationTableColumns.ReservationNumber)),
	}

	var result []exmodels.ReservationWithRelated
	err := models.Reservations(mods...).BindG(ctx, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to read reservation with related: %w", err)
	}

	return &result, nil
}

// ReadHandleTimes 店舗に紐づく、顧客＆メニューごとの施術時間一覧を取得する
func (r reservationRepository) ReadHandleTimes(ctx context.Context, storeId int) (*[]exmodels.HandleTime, error) {
	// SELECTするカラム
	selectCols := []string{
		models.ReservationTableColumns.CustomerID,
		models.ReservationTableColumns.StoreID,
		models.ReservationMenuTableColumns.MenuID,
		fmt.Sprintf("ROUND(AVG(TIMESTAMPDIFF(MINUTE, %s, %s))) as time",
			models.ReservationTableColumns.ServiceStartDatetime,
			models.ReservationTableColumns.ServiceEndDatetime,
		),
	}

	// QueryModの生成
	mods := []qm.QueryMod{
		qm.Select(strings.Join(selectCols, ",")),
		qm.InnerJoin(fmt.Sprintf("%s on %s = %s",
			models.TableNames.ReservationMenus,
			models.ReservationTableColumns.ReservationID,
			models.ReservationMenuTableColumns.ReservationID,
		)),
		qm.Where(fmt.Sprintf("%s = ?", models.ReservationTableColumns.StoreID), storeId),
		qm.Where(fmt.Sprintf("%s is not null", models.ReservationTableColumns.ServiceStartDatetime)),
		qm.Where(fmt.Sprintf("%s is not null", models.ReservationTableColumns.ServiceEndDatetime)),
		qm.GroupBy(fmt.Sprintf("%s, %s, %s",
			models.ReservationTableColumns.CustomerID,
			models.ReservationTableColumns.StoreID,
			models.ReservationMenuTableColumns.MenuID,
		)),
	}

	var result []exmodels.HandleTime
	err := models.Reservations(mods...).BindG(ctx, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to read handle time: %w", err)
	}

	return &result, nil
}
