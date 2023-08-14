package usecaseinputs

type CreateActiveStaffInput struct {
	StoreID int
	StaffID int
}

type UpdateActiveStaffInput struct {
	StoreID int
	StaffID int
	Order   int
}
type RemoveActiveStaffInput struct {
	StoreID int
	StaffID int
}
