package usecase

import (
	"app/api/errpkg"
	"app/pkg/constant"
	"app/pkg/core/repository"
	"app/pkg/enum"
	"app/pkg/models"
	"app/pkg/usecaseinputs"
	"app/pkg/usecaseoutputs"
	"context"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"time"
)

type ReservationUsecase interface {
	FetchTodayReservations(ctx context.Context, storeId int) (*[]usecaseoutputs.Reservation, error)
	FetchLineEndWaitTime(ctx context.Context, storeId int) (*usecaseoutputs.WaitTime, error)
	FetchMyWaitTime(ctx context.Context, storeId int, encryptedText string) (*usecaseoutputs.WaitTime, error)
	CreateReservation(ctx context.Context, input *usecaseinputs.CreateReservationInput) (*usecaseoutputs.CreateReservation, error)
	UpdateReservation(ctx context.Context, input *usecaseinputs.UpdateReservationInput) error
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

func (u reservationUsecase) FetchTodayReservations(ctx context.Context, storeId int) (*[]usecaseoutputs.Reservation, error) {
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
		return nil, err
	}

	for _, t := range *result {
		// 予約を特定する暗号化した文字列を生成する
		encryptedText, err := u.encryptReservation(t.ReservationID, t.StoreID, t.ReservedDatetime)
		if err != nil {
			return nil, &errpkg.UnexpectedError{
				InternalError: err,
				Operation:     "encryptReservation",
			}
		}
		reservation := usecaseoutputs.Reservation{
			ReservationID:        t.ReservationID,
			CustomerID:           t.CustomerID,
			Name:                 t.Name,
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
			MenuID:               t.MenuID,
			MenuName:             t.MenuName,
			Price:                t.Price,
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
		return nil, &errpkg.UnexpectedError{
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
		return nil, &errpkg.UnexpectedError{
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
		return nil, &errpkg.UnexpectedError{
			InternalError: err,
			Operation:     "Commit Tx",
		}
	}

	// 予約追加の通知をブロードキャストする
	u.broadcastNewReservation(2, "reservation created") // TODO storeIdハードコーディングを修正する

	// 予約を特定する暗号化した文字列を生成する
	encryptedText, err := u.encryptReservation(reservation.ReservationID, reservation.StoreID, reservation.ReservedDatetime)
	if err != nil {
		return nil, &errpkg.UnexpectedError{
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

func (u reservationUsecase) UpdateReservation(ctx context.Context, input *usecaseinputs.UpdateReservationInput) error {
	// トランザクション開始する
	tx, err := boil.BeginTx(ctx, nil)
	if err != nil {
		return &errpkg.UnexpectedError{
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

	// ステータスの値を取得する
	statusVal, exist := enum.ReservationStatusValues[input.Status]
	if !exist {
		return &errpkg.UnexpectedError{
			InternalError: err,
			Operation:     "Status not found",
		}
	}

	// 更新対象の予約を取得する
	uRes, err := u.reservationRepository.ReadReservation(ctxWithTx, input.ReservationID)
	if err != nil {
		return err
	}

	// 予約を更新する
	uRes.Status = int(statusVal)
	switch statusVal {
	case enum.Waiting: // 案内待ちに更新する場合の処理
		//...
	case enum.Pending: // 保留に更新する場合の処理
		// 保留開始時間を設定する
		uRes.HoldStartDatetime = null.TimeFrom(time.Now())
	case enum.InProgress: // 案内中に更新する場合の処理
		// 案内開始時間を設定する
		uRes.ServiceStartDatetime = null.TimeFrom(time.Now())
		// スタッフIDを設定する
		uRes.StaffID = input.StaffID
	case enum.Canceled: // キャンセルに更新する場合の処理
		// キャンセルタイプを設定する
		uRes.CancelType = null.Int{}
	case enum.Done: // 案内済みに更新する場合の処理
		// 案内終了時間を設定する
		uRes.ServiceEndDatetime = null.TimeFrom(time.Now())
	}
	// TODO
	uRes.ArrivalFlag = false

	_, err = u.reservationRepository.UpdateReservation(ctxWithTx, uRes)
	if err != nil {
		return err
	}

	// トランザクションをコミットする
	err = tx.Commit()
	if err != nil {
		return &errpkg.UnexpectedError{
			InternalError: err,
			Operation:     "Commit Tx",
		}
	}

	return nil
}
