package usecase

import (
	"app/api/errors"
	"app/pkg/constant"
	"app/pkg/core/repository"
	"app/pkg/enum"
	"app/pkg/models"
	"app/pkg/usecaseinputs"
	"app/pkg/usecaseoutputs"
	"context"
	"fmt"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"time"
)

type ReservationUsecase interface {
	FetchAllReservations(ctx context.Context, storeId int) (*[]usecaseoutputs.Reservation, error)
	FetchLineEndWaitTime(ctx context.Context, storeId int) (*usecaseoutputs.WaitTime, error)
	FetchMyWaitTime(ctx context.Context, storeId int, encryptedText string) (*usecaseoutputs.WaitTime, error)
	CreateReservation(ctx context.Context, input *usecaseinputs.CreateReservationInput) (*usecaseoutputs.CreateReservation, error)
}

type reservationUsecase struct {
	reservationRepository     repository.ReservationRepository
	reservationMenuRepository repository.ReservationMenuRepository
	activeStaffRepository     repository.ActiveStaffRepository
}

func NewReservationUsecase(
	rr repository.ReservationRepository,
	rmr repository.ReservationMenuRepository,
	asr repository.ActiveStaffRepository,
) ReservationUsecase {
	return &reservationUsecase{
		reservationRepository:     rr,
		reservationMenuRepository: rmr,
		activeStaffRepository:     asr,
	}
}

func (u reservationUsecase) FetchAllReservations(ctx context.Context, storeId int) (*[]usecaseoutputs.Reservation, error) {
	reservations := make([]usecaseoutputs.Reservation, 0)
	result, err := u.reservationRepository.ReadTodayReservations(
		ctx,
		storeId,
		enum.Waiting,
		enum.InProgress,
		enum.Done,
		enum.Pending,
		enum.Canceled,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to Fetch all reservations: %w", err)
	}

	for _, t := range *result {
		// 予約を特定する暗号化した文字列を生成する
		encryptedText, err := u.encryptReservation(t.ReservationID, t.StoreID, t.ReservedDatetime)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt reservation: %w", err)
		}
		reservation := usecaseoutputs.Reservation{
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
			Content:              encryptedText,
		}
		reservations = append(reservations, reservation)
	}

	return &reservations, nil
}

func (u reservationUsecase) FetchLineEndWaitTime(ctx context.Context, storeId int) (*usecaseoutputs.WaitTime, error) {
	// 施術目安時間更新後の本日の予約一覧を取得する
	reservations, err := u.fetchReservationsWithUpdateHandleTimes(ctx, storeId)
	if err != nil {
		return nil, err
	}

	// 次の予約番号
	nextReservationNumber, err := u.getNextReservationNumber(ctx, storeId)
	if err != nil {
		return nil, err
	}
	// 次の順番
	position := 1
	// 現在の待ち時間
	waitTime := 0

	// 施術中スタッフ一覧スタッフの対応可能時間を取得する
	activeStaffs, err := u.activeStaffRepository.ReadActiveStaffs(ctx, storeId)
	if err != nil {
		return nil, err
	}
	staffAvailableTimes := u.getStaffAvailableTimes(activeStaffs, reservations)

	// 案内待ちの予約数＋1番目の数を取得する
	position = u.getNextPosition(reservations)
	// 待ち時間を取得する
	waitTime = calcWaitTime(*reservations, *staffAvailableTimes, len(*reservations))

	return &usecaseoutputs.WaitTime{
		ReservationNumber: nextReservationNumber,
		Position:          position,
		Time:              waitTime,
	}, nil
}

func (u reservationUsecase) FetchMyWaitTime(ctx context.Context, storeId int, encryptedText string) (*usecaseoutputs.WaitTime, error) {
	// 暗号化文字列を復号化する
	reservationInfo, err := u.decryptReservation(encryptedText)
	if err != nil {
		return nil, &errors.UnexpectedError{
			InternalError: err,
			Operation:     "decryptReservation",
		}
	}

	// 施術目安時間更新後の本日の予約一覧を取得する
	reservations, err := u.fetchReservationsWithUpdateHandleTimes(ctx, storeId)
	if err != nil {
		return nil, err
	}

	// 順番、予約番号
	var position, reservationNumber int
	// 予約IDに紐づく予約が何番目か特定する
	for i, rs := range *reservations {
		if rs.ReservationID == reservationInfo.ReservationID {
			position = i + 1
			reservationNumber = rs.ReservationNumber
		}
	}

	// 施術中スタッフ一覧スタッフの対応可能時間を取得する
	activeStaffs, _ := u.activeStaffRepository.ReadActiveStaffs(ctx, storeId)
	staffAvailableTimes := u.getStaffAvailableTimes(activeStaffs, reservations)

	// 待ち時間を取得する
	waitTime := calcWaitTime(*reservations, *staffAvailableTimes, position)

	return &usecaseoutputs.WaitTime{
		ReservationNumber: reservationNumber,
		Position:          position,
		Time:              waitTime,
	}, nil
}

func (u reservationUsecase) CreateReservation(ctx context.Context, input *usecaseinputs.CreateReservationInput) (*usecaseoutputs.CreateReservation, error) {
	// トランザクション開始する
	tx, err := boil.BeginTx(ctx, nil)
	if err != nil {
		return nil, &errors.UnexpectedError{
			InternalError: err,
			Operation:     "Begin Tx",
		}
	}
	// トランザクションをcontextに紐づける
	ctxWithTx := context.WithValue(ctx, constant.ContextExecutorKey, tx)

	// エラーの場合ロールバックする
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// 次の予約番号を取得する
	nextReservationNumber, err := u.getNextReservationNumber(ctxWithTx, input.StoreID)
	if err != nil {
		return nil, err
	}

	// 予約を登録する
	cRes := &models.Reservation{
		CustomerID:        input.CustomerID,
		StoreID:           input.StoreID,
		ReservationNumber: nextReservationNumber,
		ReservedDatetime:  time.Now(),
		Status:            int(enum.Waiting),
	}
	reservation, err := u.reservationRepository.InsertReservation(ctxWithTx, cRes)
	if err != nil {
		return nil, err
	}

	// 予約メニューを登録する
	cResMn := &models.ReservationMenu{
		ReservationID: cRes.ReservationID,
		StoreID:       input.StoreID,
		MenuID:        input.MenuID,
	}
	err = u.reservationMenuRepository.InsertReservationMenu(ctx, cResMn)
	if err != nil {
		return nil, err
	}

	// トランザクションをコミットする
	err = tx.Commit()
	if err != nil {
		return nil, &errors.UnexpectedError{
			InternalError: err,
			Operation:     "Commit Tx",
		}
	}

	// 予約追加の通知をブロードキャストする
	u.broadcastNewReservation(2, "reservation created") // TODO storeIdハードコーディングを修正する

	// 予約を特定する暗号化した文字列を生成する
	encryptedText, err := u.encryptReservation(reservation.ReservationID, reservation.StoreID, reservation.ReservedDatetime)
	if err != nil {
		return nil, &errors.UnexpectedError{
			InternalError: err,
			Operation:     "encryptReservation",
		}
	}

	output := &usecaseoutputs.CreateReservation{
		ReservationNumber: reservation.ReservationNumber,
		Content:           encryptedText,
	}

	return output, nil
}
