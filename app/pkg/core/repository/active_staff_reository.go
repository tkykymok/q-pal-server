package repository

import (
	"app/api/errpkg"
	"app/pkg/models"
	"context"
	"fmt"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type ActiveStaffRepository interface {
	ReadActiveStaff(ctx context.Context, activeStaff *models.ActiveStaff) (*models.ActiveStaff, error)

	ReadActiveStaffs(ctx context.Context, storeId int) (models.ActiveStaffSlice, error)

	InsertActiveStaff(ctx context.Context, activeStaff *models.ActiveStaff) error

	UpdateActiveStaff(ctx context.Context, activeStaff *models.ActiveStaff) error

	DeleteActiveStaff(ctx context.Context, activeStaff *models.ActiveStaff) error

	DeleteActiveStaffs(ctx context.Context, storeId int) error
}

type activeStaffRepository struct {
}

func NewActiveStaffRepo() ActiveStaffRepository {
	return &activeStaffRepository{}
}

func (r activeStaffRepository) ReadActiveStaff(ctx context.Context, activeStaff *models.ActiveStaff) (*models.ActiveStaff, error) {
	// QueryModの生成
	mods := []qm.QueryMod{
		qm.Where(fmt.Sprintf("%s = ?", models.ActiveStaffTableColumns.StoreID), activeStaff.StoreID),
		qm.Where(fmt.Sprintf("%s = ?", models.ActiveStaffTableColumns.StaffID), activeStaff.StaffID),
	}

	result, err := models.ActiveStaffs(mods...).OneG(ctx)
	if err != nil {
		return nil, &errpkg.DatabaseError{
			InternalError: err,
			Operation:     "ReadActiveStaff",
		}
	}
	return result, nil
}

func (r activeStaffRepository) ReadActiveStaffs(ctx context.Context, storeId int) (models.ActiveStaffSlice, error) {
	// QueryModの生成
	mods := []qm.QueryMod{
		qm.Where(fmt.Sprintf("%s = ?", models.ActiveStaffTableColumns.StoreID), storeId),
		qm.OrderBy(fmt.Sprintf("%s asc", models.ActiveStaffTableColumns.Order)),
	}

	result, err := models.ActiveStaffs(mods...).AllG(ctx)
	if err != nil {
		return nil, &errpkg.DatabaseError{
			InternalError: err,
			Operation:     "ReadActiveStaffs",
		}
	}
	return result, nil
}

func (r activeStaffRepository) InsertActiveStaff(ctx context.Context, activeStaff *models.ActiveStaff) error {
	err := activeStaff.InsertG(ctx, boil.Infer())
	if err != nil {
		return &errpkg.DatabaseError{
			InternalError: err,
			Operation:     "InsertActiveStaff",
		}
	}
	return nil
}

func (r activeStaffRepository) UpdateActiveStaff(ctx context.Context, activeStaff *models.ActiveStaff) error {
	_, err := activeStaff.UpdateG(ctx, boil.Infer())
	if err != nil {
		return &errpkg.DatabaseError{
			InternalError: err,
			Operation:     "UpdateActiveStaff",
		}
	}
	return nil
}

func (r activeStaffRepository) DeleteActiveStaff(ctx context.Context, activeStaff *models.ActiveStaff) error {
	_, err := activeStaff.DeleteG(ctx)
	if err != nil {
		return &errpkg.DatabaseError{
			InternalError: err,
			Operation:     "DeleteActiveStaff",
		}
	}
	return nil
}

func (r activeStaffRepository) DeleteActiveStaffs(ctx context.Context, storeId int) error {
	// QueryModの生成
	mods := []qm.QueryMod{
		qm.Where(fmt.Sprintf("%s = ?", models.ActiveStaffTableColumns.StoreID), storeId),
	}

	_, err := models.ActiveStaffs(mods...).DeleteAllG(ctx)

	if err != nil {
		return &errpkg.DatabaseError{
			InternalError: err,
			Operation:     "DeleteActiveStaffs",
		}
	}
	return nil
}
