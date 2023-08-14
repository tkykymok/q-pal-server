package usecase

import (
	"app/api/errpkg"
	"app/pkg/constant"
	"app/pkg/core/repository"
	"app/pkg/models"
	"app/pkg/usecaseinputs"
	"app/pkg/usecaseoutputs"
	"context"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"time"
)

type StaffUsecase interface {
	FetchStaffs(ctx context.Context, storeId int) (*[]usecaseoutputs.Staff, error)

	CreateActiveStaff(ctx context.Context, input usecaseinputs.CreateActiveStaffInput) error

	RemoveActiveStaff(ctx context.Context, input usecaseinputs.RemoveActiveStaffInput) error
}

type staffUsecase struct {
	staffRepository       repository.StaffRepository
	activeStaffRepository repository.ActiveStaffRepository
}

func NewStaffUsecase(
	sr repository.StaffRepository,
	asr repository.ActiveStaffRepository,
) StaffUsecase {
	return &staffUsecase{
		staffRepository:       sr,
		activeStaffRepository: asr,
	}
}

func (u staffUsecase) FetchStaffs(ctx context.Context, storeId int) (*[]usecaseoutputs.Staff, error) {
	staffs := make([]usecaseoutputs.Staff, 0)

	result, err := u.staffRepository.ReadStaffs(ctx, storeId)
	if err != nil {
		return nil, err
	}

	for _, t := range *result {
		reservation := usecaseoutputs.Staff{
			StoreID:            t.StoreID,
			StaffID:            t.StaffID,
			Name:               t.Name,
			Order:              t.Order,
			BreakStartDatetime: t.BreakStartDatetime,
			BreakEndDatetime:   t.BreakEndDatetime,
			ActiveFlag:         t.ActiveFlag,
		}
		staffs = append(staffs, reservation)
	}

	return &staffs, nil
}

func (u staffUsecase) CreateActiveStaff(ctx context.Context, input usecaseinputs.CreateActiveStaffInput) error {
	// アクティブスタッフ一覧を取得する
	activeStaffs, err := u.activeStaffRepository.ReadActiveStaffs(ctx, input.StoreID)
	if err != nil {
		return err
	}

	// 追加するアクティブスタッフの構造体を生成する
	activeStaff := &models.ActiveStaff{
		StaffID: input.StaffID,
		StoreID: input.StoreID,
		Order:   len(activeStaffs) + 1,
	}

	// アクティブスタッフを登録する
	err = u.activeStaffRepository.InsertActiveStaff(ctx, activeStaff)
	if err != nil {
		return err
	}

	return nil
}

func (u staffUsecase) RemoveActiveStaff(ctx context.Context, input usecaseinputs.RemoveActiveStaffInput) error {
	// 削除するアクティブスタッフの構造体を生成する
	activeStaff := &models.ActiveStaff{
		StaffID: input.StaffID,
		StoreID: input.StoreID,
	}

	// 削除対象アクティブスタッフが存在するかチェック
	result, err := u.activeStaffRepository.ReadActiveStaff(ctx, activeStaff)
	if err != nil {
		return &errpkg.UnexpectedError{
			InternalError: err,
			Operation:     "Check Staff Existence",
		}
	}

	// トランザクション開始する
	tx, err := boil.BeginTx(ctx, nil)
	if err != nil {
		return &errpkg.UnexpectedError{
			InternalError: err,
			Operation:     "Begin Tx",
		}
	}
	// エラーの場合ロールバックする
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()
	// トランザクションをcontextに紐づける
	ctxWithTx := context.WithValue(ctx, constant.ContextExecutorKey, tx)

	// アクティブスタッフを削除する
	err = u.activeStaffRepository.DeleteActiveStaff(ctxWithTx, result)
	if err != nil {
		return err
	}

	// アクティブスタッフ一覧を取得する
	activeStaffs, err := u.activeStaffRepository.ReadActiveStaffs(ctx, input.StoreID)
	if err != nil {
		return err
	}

	// 表示順を更新する
	for i, as := range activeStaffs {
		as.Order = i + 1
		as.UpdatedAt = time.Now()
		err = u.activeStaffRepository.UpdateActiveStaff(ctxWithTx, as)
		if err != nil {
			return err
		}
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
