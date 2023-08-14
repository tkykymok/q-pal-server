package repository

import (
	"app/api/errpkg"
	"app/pkg/exmodels"
	"app/pkg/models"
	"context"
	"fmt"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"strings"
)

type StaffRepository interface {
	ReadStaffs(ctx context.Context, storeId int) (*[]exmodels.StaffWithRelated, error)
}

type staffRepository struct {
}

func NewStaffRepo() StaffRepository {
	return &staffRepository{}
}

func (s staffRepository) ReadStaffs(ctx context.Context, storeId int) (*[]exmodels.StaffWithRelated, error) {
	// SELECTするカラム
	selectCols := []string{
		models.StoreStaffTableColumns.StoreID,
		models.StaffTableColumns.StaffID,
		models.StaffTableColumns.Name,
		models.ActiveStaffTableColumns.Order,
		models.ActiveStaffTableColumns.BreakStartDatetime,
		models.ActiveStaffTableColumns.BreakEndDatetime,
		fmt.Sprintf("case when %s is null then false else true end as active_flag",
			models.ActiveStaffTableColumns.Order),
	}

	// QueryModの生成
	mods := []qm.QueryMod{
		qm.Select(strings.Join(selectCols, ",")),
		qm.InnerJoin(fmt.Sprintf("%s on %s = %s",
			models.TableNames.StoreStaffs,
			models.StaffTableColumns.StaffID,
			models.StoreStaffTableColumns.StaffID,
		)),
		qm.LeftOuterJoin(fmt.Sprintf("%s on %s = %s ",
			models.TableNames.ActiveStaffs,
			models.StoreStaffTableColumns.StaffID,
			models.ActiveStaffTableColumns.StaffID,
		)),
		qm.Where(fmt.Sprintf("%s = ?", models.StoreStaffTableColumns.StoreID), storeId),
	}

	var result []exmodels.StaffWithRelated
	err := models.Staffs(mods...).BindG(ctx, &result)
	if err != nil {
		return nil, &errpkg.DatabaseError{
			InternalError: err,
			Operation:     "ReadStaffs",
		}
	}

	return &result, nil
}
