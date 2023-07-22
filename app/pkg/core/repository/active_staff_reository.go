package repository

import (
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

	return models.ActiveStaffs(mods...).AllG(ctx)
}
