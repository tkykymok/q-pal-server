package repository

import (
	"app/api/errors"
	"app/pkg/constant"
	"app/pkg/enum"
	"app/pkg/exmodels"
	"app/pkg/models"
	"app/pkg/utils"
	"context"
	"fmt"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"strings"
	"time"
)

type ReservationRepository interface {
	ReadTodayReservations(ctx context.Context, storeId int, status ...enum.ReservationStatus) (*[]exmodels.ReservationWithRelated, error)

	ReadLatestReservation(ctx context.Context, storeId int) (models.ReservationSlice, error)

	ReadHandleTimes(ctx context.Context, storeId int) (*[]exmodels.HandleTime, error)

	InsertReservation(ctx context.Context, reservation *models.Reservation) (*models.Reservation, error)
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
		qm.Where(fmt.Sprintf("%s >= ?", models.ReservationTableColumns.ReservedDatetime), utils.GetStartOfDay()),
		qm.Where(fmt.Sprintf("%s <= ?", models.ReservationTableColumns.ReservedDatetime), utils.GetEndOfDay()),
		qm.WhereIn(fmt.Sprintf("%s in ?", models.ReservationTableColumns.Status), statusInts...),
		qm.OrderBy(fmt.Sprintf("%s", models.ReservationTableColumns.ReservationNumber)),
	}

	var result []exmodels.ReservationWithRelated
	err := models.Reservations(mods...).BindG(ctx, &result)
	if err != nil {
		return nil, &errors.DatabaseError{
			InternalError: err,
			Operation:     "ReadTodayReservations",
		}
	}

	return &result, nil
}

func (r reservationRepository) ReadLatestReservation(ctx context.Context, storeId int) (models.ReservationSlice, error) {
	// QueryModの生成
	mods := []qm.QueryMod{
		qm.Where(fmt.Sprintf("%s = ?", models.ReservationTableColumns.StoreID), storeId),
		qm.Where(fmt.Sprintf("%s >= ?", models.ReservationTableColumns.ReservedDatetime), utils.GetStartOfDay()),
		qm.Where(fmt.Sprintf("%s <= ?", models.ReservationTableColumns.ReservedDatetime), utils.GetEndOfDay()),
		qm.OrderBy(fmt.Sprintf("%s desc", models.ReservationTableColumns.ReservationNumber)),
		qm.Limit(1),
	}

	// トランザクションが存在する場合のみ、タイムアウトを設定
	if ctx.Value(constant.ContextExecutorKey) != nil {
		mods = append(mods, qm.For("update"))
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
	}

	result, err := models.Reservations(mods...).AllG(ctx)
	if err != nil {
		return nil, &errors.DatabaseError{
			InternalError: err,
			Operation:     "ReadLatestReservation",
		}
	}

	return result, nil
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
		return nil, &errors.DatabaseError{
			InternalError: err,
			Operation:     "ReadHandleTimes",
		}
	}

	return &result, nil
}

func (r reservationRepository) InsertReservation(ctx context.Context, reservation *models.Reservation) (*models.Reservation, error) {
	err := reservation.InsertG(ctx, boil.Infer())
	if err != nil {
		return nil, &errors.DatabaseError{
			InternalError: err,
			Operation:     "InsertReservation",
		}
	}

	return reservation, nil
}
