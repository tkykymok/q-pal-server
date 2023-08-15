package usecaseinputs

import "github.com/volatiletech/null/v8"

type CreateActiveStaffInput struct {
	StoreID int
	StaffID int
	Order   null.Int
}

type UpdateActiveStaffData struct {
	StaffID int
	Order   int
}

type UpdateActiveStaffInput struct {
	StoreId int
	Data    []UpdateActiveStaffData
}

type RemoveActiveStaffInput struct {
	StoreID int
	StaffID int
}
