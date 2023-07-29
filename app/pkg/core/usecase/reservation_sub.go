package usecase

import (
	"app/api/presenter"
	"app/pkg/broadcast"
	"app/pkg/enum"
	"app/pkg/exmodels"
	"app/pkg/models"
	"container/heap"
	"context"
	"fmt"
	"github.com/volatiletech/null/v8"
	"time"
)

type StaffAvailableTime struct {
	StaffID    int // スタッフID
	FinishTime int // 対応終了時間(現在時刻から何分後)
	BreakStart int // 休憩開始時間(現在時刻から何分後)
	BreakEnd   int // 休憩終了時間(現在時刻から何分後)
}

type PriorityQueue []*StaffAvailableTime

func (pq *PriorityQueue) Len() int { return len(*pq) }

func (pq *PriorityQueue) Less(i, j int) bool {
	return (*pq)[i].FinishTime < (*pq)[j].FinishTime
}

func (pq *PriorityQueue) Swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*StaffAvailableTime))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[0 : n-1]
	return x
}

// updateTimes  予約一覧に対する想定施術時間を更新する
func (u reservationUsecase) updateTimes(reservations *[]exmodels.ReservationWithRelated, handleTimes *[]exmodels.HandleTime) {
	// 予約IDをキーに施術時間をMapに格納する
	timeMap := make(map[string]int)
	for _, ht := range *handleTimes {
		key := fmt.Sprintf("%d-%d-%d", ht.CustomerID, ht.StoreID, ht.MenuID)
		timeMap[key] = ht.Time
	}

	// 予約一覧に対する施術時間を更新する
	for i, r := range *reservations {
		key := fmt.Sprintf("%d-%d-%d", r.CustomerID, r.StoreID, r.MenuID)
		if handleTime, ok := timeMap[key]; ok {
			(*reservations)[i].Time = handleTime
		}
	}
}

// getStaffAvailableTimes スタッフの現在対応状況を取得する(何分後に終了するか & 休憩予定)
func (u reservationUsecase) getStaffAvailableTimes(activeStaffs models.ActiveStaffSlice, reservations *[]exmodels.ReservationWithRelated) *[]StaffAvailableTime {
	staffAvailableTimes := make([]StaffAvailableTime, 0, len(activeStaffs))
	for _, as := range activeStaffs {
		finishTime := 0
		for _, r := range *reservations {
			if r.StaffID.Valid && r.StaffID.Int == as.StaffID {
				passedTime := int(time.Now().Sub(r.ServiceStartDatetime.Time).Minutes())
				if r.Time-passedTime > 0 {
					finishTime = r.Time - passedTime
				}
			}
		}
		staffAvailableTime := &StaffAvailableTime{
			StaffID:    as.StaffID,
			FinishTime: finishTime,
			BreakStart: calcDiffTime(as.BreakStartDatetime),
			BreakEnd:   calcDiffTime(as.BreakEndDatetime),
		}
		staffAvailableTimes = append(staffAvailableTimes, *staffAvailableTime) // スライスに要素を追加
	}
	return &staffAvailableTimes
}

// calcWaitTime 待ち時間を計算する
func calcWaitTime(reservations []exmodels.ReservationWithRelated, staffAvailableTimes []StaffAvailableTime, position int) int {
	// Create a copy of staffs slices
	copiedStaffs := make([]*StaffAvailableTime, len(staffAvailableTimes))
	for i, staff := range staffAvailableTimes {
		copiedStaff := staff
		copiedStaffs[i] = &copiedStaff
	}

	staffAvailableTimesQueue := &PriorityQueue{}
	heap.Init(staffAvailableTimesQueue)
	for _, staff := range copiedStaffs {
		heap.Push(staffAvailableTimesQueue, staff)
	}

	// 待ち時間
	waitTime := 0
	for i, rs := range reservations {
		staff := heap.Pop(staffAvailableTimesQueue).(*StaffAvailableTime)
		if i+1 == position {
			// ループ対象の予約が案内中以外の場合
			if rs.Status != int(enum.InProgress) {
				// 待ち時間を設定する
				waitTime = staff.FinishTime
			}
			break
		}

		finishTime := staff.FinishTime + reservations[i].Time
		if !(staff.BreakStart == 0 && staff.BreakEnd == 0) {
			if staff.BreakStart <= staff.FinishTime && staff.FinishTime < staff.BreakEnd {
				finishTime = staff.BreakEnd
			}
		}

		staff.FinishTime = finishTime
		heap.Push(staffAvailableTimesQueue, staff)
	}

	roundedWaitTime := int((float64(waitTime)/5)+0.5) * 5
	return roundedWaitTime
}

// 現在時刻との差分(分単位)を計算する
func calcDiffTime(t null.Time) int {
	if t.Valid {
		diff := int(t.Time.Sub(time.Now()).Minutes())
		if diff > 0 {
			return diff
		}
	}
	return 0
}

// 店舗に紐づく最新予約番号を取得する
func (u reservationUsecase) getNextReservationNumber(ctx context.Context, storeId int) (int, error) {
	reservationNumber := 1 // 初期値
	// 全ステータスの予約一覧を取得
	reservations, err := u.reservationRepository.ReadLatestReservation(ctx, storeId)
	if err != nil {
		return 0, fmt.Errorf("failed to get next reservation number: %w", err)
	}
	// 検索結果が存在する場合
	if reservations != nil {
		return (reservations)[len(reservations)-1].ReservationNumber + 1, nil
	}
	return reservationNumber, nil
}

// 案内待ちの予約数＋1番目の数を返す
func (u reservationUsecase) getNextPosition(reservations *[]exmodels.ReservationWithRelated) int {
	count := 1
	for _, rs := range *reservations {
		if rs.Status == int(enum.Waiting) {
			count++
		}
	}
	return count
}

func (u reservationUsecase) broadcastNewReservation(storeId int, message string) {
	messsage := presenter.ReservationMessage{
		Message: message,
	}

	broadcast.ReservationClient.SendNewReservation(storeId, messsage)
}
