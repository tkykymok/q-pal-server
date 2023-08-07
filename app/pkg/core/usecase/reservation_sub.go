package usecase

import (
	"app/api/presenter"
	"app/pkg/broadcast"
	"app/pkg/enum"
	"app/pkg/exmodels"
	"app/pkg/models"
	"app/pkg/usecaseoutputs"
	"container/heap"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null/v8"
	"io"
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

// 施術目安時間更新後の本日の予約一覧を取得する
func (u reservationUsecase) fetchReservationsWithUpdateHandleTimes(ctx context.Context, storeId int) (*[]exmodels.ReservationWithRelated, error) {
	// 本日の予約一覧を取得
	reservations, err := u.reservationRepository.ReadTodayReservations(
		ctx,
		storeId,
		enum.Waiting,    // 未案内
		enum.InProgress, // 案内中
	)
	if err != nil {
		return nil, err
	}

	// 顧客ごとの過去履歴に紐づく施術時間一覧を取得する
	handleTimes, err := u.reservationRepository.ReadHandleTimes(ctx, storeId)
	if err != nil {
		return nil, err
	}

	// 予約一覧に対する施術時間を更新する
	u.updateTimes(reservations, handleTimes)

	return reservations, nil
}

// 店舗に紐づく最新予約番号を取得する
func (u reservationUsecase) getNextReservationNumber(ctx context.Context, storeId int) (int, error) {
	reservationNumber := 1 // 初期値
	// 全ステータスの予約一覧を取得
	reservations, err := u.reservationRepository.ReadLatestReservation(ctx, storeId)
	if err != nil {
		return 0, err
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
	reservationMessage := presenter.ReservationMessage{
		Message: message,
	}

	broadcast.ReservationClient.SendNewReservation(storeId, reservationMessage)
}

func (u reservationUsecase) encryptReservation(reservationId int, storeId int, reservedDatetime time.Time) (string, error) {
	// 暗号化キーを設定します。
	encryptKey := "pass1234pass1234"

	// 暗号化するためのデータを構造体にセットします。
	reservationKey := usecaseoutputs.ReservationIdentifyKey{
		ReservationID:    reservationId,
		StoreID:          storeId,
		ReservedDatetime: reservedDatetime,
	}

	// 構造体をJSON形式に変換します。
	jsonData, err := json.Marshal(reservationKey)
	if err != nil {
		return "", err
	}

	// AES暗号化アルゴリズムのための暗号化キーを生成します。
	block, err := aes.NewCipher([]byte(encryptKey))
	if err != nil {
		return "", err
	}

	// GCM（Galois/Counter Mode）というAESのモードを生成します。これにより、暗号化と同時に改ざん検知が可能になります。
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// ランダムなnonce（Number used once）を生成します。これは暗号化の初期ベクトルとして使用され、同じデータでも異なる暗号文を生成するために使われます。
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// 暗号化を行い、nonceを暗号文の先頭に追加します。これにより、復号化時に同じnonceを使用できます。
	ciphertext := gcm.Seal(nonce, nonce, jsonData, nil)

	// 暗号文をBase64形式の文字列に変換します。これにより、バイナリデータをテキスト形式で安全に扱うことができます。
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (u reservationUsecase) decryptReservation(ciphertextStr string) (*usecaseoutputs.ReservationIdentifyKey, error) {
	if len(ciphertextStr) == 0 {
		return nil, errors.New("Invalid parameter")
	}

	// 暗号化キーを設定します。
	encryptKey := "pass1234pass1234"

	// Base64形式の文字列をバイナリデータに戻します。
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextStr)
	if err != nil {
		return nil, err
	}

	// AES暗号化アルゴリズムのための暗号化キーを生成します。
	block, err := aes.NewCipher([]byte(encryptKey))
	if err != nil {
		return nil, err
	}

	// GCM（Galois/Counter Mode）というAESのモードを生成します。これにより、暗号化と同時に改ざん検知が可能になります。
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// nonceのサイズを取得します。
	nonceSize := gcm.NonceSize()

	// 暗号文からnonceと実際の暗号文を分離します。
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// 暗号文を復号化します。
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	// 復号化したデータを構造体に戻します。
	var reservationKey usecaseoutputs.ReservationIdentifyKey
	err = json.Unmarshal(plaintext, &reservationKey)
	if err != nil {
		return nil, err
	}

	return &reservationKey, nil
}
