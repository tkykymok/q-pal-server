package usecase

import (
	"app/pkg/core/repository"
	"app/pkg/models"
	"app/pkg/usecaseinputs"
	"app/pkg/usecaseoutputs"
	"context"
)

type StaffUsecase interface {
	FetchStaffs(ctx context.Context, storeId int) (*[]usecaseoutputs.Staff, error)

	CreateActiveStaff(ctx context.Context, input usecaseinputs.CreateActiveStaffInput) error
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
	// アクティブなスタッフ一覧を取得する
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
