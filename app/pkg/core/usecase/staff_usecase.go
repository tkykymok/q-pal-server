package usecase

import (
	"app/pkg/core/repository"
	"app/pkg/usecaseoutputs"
	"context"
)

type StaffUsecase interface {
	FetchStaffs(ctx context.Context, storeId int) (*[]usecaseoutputs.Staff, error)
}

type staffUsecase struct {
	staffRepository repository.StaffRepository
}

func NewStaffUsecase(
	sr repository.StaffRepository,
) StaffUsecase {
	return &staffUsecase{
		staffRepository: sr,
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
