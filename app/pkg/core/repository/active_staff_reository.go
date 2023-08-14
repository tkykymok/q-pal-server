package repository

import (
	"app/api/errors"
	"app/pkg/models"
	"context"
	"fmt"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type ActiveStaffRepository interface {
	ReadActiveStaffs(ctx context.Context, storeId int) (models.ActiveStaffSlice, error)

	InsertActiveStaff(ctx context.Context, activeStaff *models.ActiveStaff) error
}

type activeStaffRepository struct {
}

func NewActiveStaffRepo() ActiveStaffRepository {
	return &activeStaffRepository{}
}

func (r activeStaffRepository) ReadActiveStaffs(ctx context.Context, storeId int) (models.ActiveStaffSlice, error) {
	// QueryModの生成
	mods := []qm.QueryMod{
		qm.Where(fmt.Sprintf("%s = ?", models.ActiveStaffTableColumns.StoreID), storeId),
	}

	result, err := models.ActiveStaffs(mods...).AllG(ctx)
	if err != nil {
		return nil, &errors.DatabaseError{
			InternalError: err,
			Operation:     "ReadActiveStaffs",
		}
	}
	return result, nil
}

func (r activeStaffRepository) InsertActiveStaff(ctx context.Context, activeStaff *models.ActiveStaff) error {
	err := activeStaff.InsertG(ctx, boil.Infer())
	if err != nil {
		return &errors.DatabaseError{
			InternalError: err,
			Operation:     "InsertActiveStaff",
		}
	}
	return nil
}
