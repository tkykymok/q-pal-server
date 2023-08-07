package repository

import (
	"app/api/errors"
	"app/pkg/models"
	"context"
	"fmt"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type ActiveStaffRepository interface {
	ReadActiveStaffs(ctx context.Context, storeId int) (models.ActiveStaffSlice, error)
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
