package requests

import "github.com/volatiletech/null/v8"

type CreateActiveStaff struct {
	StaffId int
	Order   null.Int
}

type UpdateActiveStaff struct {
	Data []struct {
		StaffId int
		Order   int
	}
}
