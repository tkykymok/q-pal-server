package usecase

import (
	"app/pkg/core/repository"
	"app/pkg/enum"
	"app/pkg/outputs"
	"context"
	"fmt"
)

type ReservationUsecase interface {
	FetchAllReservations(ctx context.Context, storeId int) (*[]outputs.Reservation, error)
	FetchLineEndWaitTime(ctx context.Context, storeId int) (*outputs.WaitTime, error)
	FetchIndividualWaitTime(ctx context.Context, storeId int, reservationId int) (*outputs.WaitTime, error)
}

type reservationUsecase struct {
	reservationRepository repository.ReservationRepository
	activeStaffRepository repository.ActiveStaffRepository
}

func NewReservationUsecase(rr repository.ReservationRepository, asr repository.ActiveStaffRepository) ReservationUsecase {
	return &reservationUsecase{
		reservationRepository: rr,
		activeStaffRepository: asr,
	}
}

func (u reservationUsecase) FetchAllReservations(ctx context.Context, storeId int) (*[]outputs.Reservation, error) {
	reservations := make([]outputs.Reservation, 0)
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
		return nil, fmt.Errorf("failed to read today reservations: %w", err)
	}

	for _, t := range *result {
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

func (u reservationUsecase) FetchLineEndWaitTime(ctx context.Context, storeId int) (*outputs.WaitTime, error) {
	// 本日の予約一覧を取得
	reservations, err := u.reservationRepository.ReadTodayReservations(
		ctx,
		storeId,
		enum.Waiting,    // 未案内
		enum.InProgress, // 案内中
	)
	if err != nil {
		return nil, fmt.Errorf("failed to read line end wai time: %w", err)
	}

	// 次の予約番号
	nextReservationNumber := u.getNextReservationNumber(ctx, storeId)
	// 次の順番
	position := 1
	// 現在の待ち時間
	waitTime := 0
	if len(*reservations) > 0 {
		// 顧客ごとの過去履歴に紐づく施術時間一覧を取得する
		handleTimes, err := u.reservationRepository.ReadHandleTimes(ctx, storeId)
		if err != nil {
			return nil, fmt.Errorf("failed to read handle times: %w", err)
		}

		// 予約一覧に対する施術時間を更新する
		u.updateTimes(reservations, handleTimes)

		// 施術中スタッフ一覧スタッフの対応可能時間を取得する
		activeStaffs, err := u.activeStaffRepository.ReadActiveStaffs(ctx, storeId)
		if err != nil {
			return nil, fmt.Errorf("failed to read active staffs: %w", err)
		}
		staffAvailableTimes := u.getStaffAvailableTimes(activeStaffs, reservations)

		position = u.getNextPosition(reservations)
		waitTime = calcWaitTime(*reservations, *staffAvailableTimes, len(*reservations))
	}

	return &outputs.WaitTime{
		ReservationNumber: nextReservationNumber,
		Position:          position,
		Time:              waitTime,
	}, nil
}

func (u reservationUsecase) FetchIndividualWaitTime(ctx context.Context, storeId int, reservationId int) (*outputs.WaitTime, error) {
	// 予約一覧を取得
	reservations, _ := u.reservationRepository.ReadTodayReservations(
		ctx,
		storeId,
		enum.Waiting,
		enum.InProgress,
	)

	// 順番
	var position int
	// 予約番号
	var reservationNumber int
	// 予約IDに紐づく予約が何番目か特定する
	for i, rs := range *reservations {
		if rs.ReservationID == reservationId {
			position = i + 1
			reservationNumber = rs.ReservationNumber
		}
	}

	// 顧客ごとの過去履歴に紐づく施術時間一覧を取得する
	handleTimes, _ := u.reservationRepository.ReadHandleTimes(ctx, storeId)

	// 予約一覧に対する施術時間を更新する
	u.updateTimes(reservations, handleTimes)

	// 施術中スタッフ一覧スタッフの対応可能時間を取得する
	activeStaffs, _ := u.activeStaffRepository.ReadActiveStaffs(ctx, storeId)
	staffAvailableTimes := u.getStaffAvailableTimes(activeStaffs, reservations)

	waitTime := calcWaitTime(*reservations, *staffAvailableTimes, position)

	return &outputs.WaitTime{
		ReservationNumber: reservationNumber,
		Position:          position,
		Time:              waitTime,
	}, nil
}
