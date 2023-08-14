package presenter

import (
	"app/pkg/usecaseoutputs"
	"app/pkg/utils"
	"github.com/volatiletech/null/v8"
)

type Staff struct {
	StoreID            int         `json:"storeId"`
	StaffID            int         `json:"staffId"`
	Name               null.String `json:"name"`
	Order              null.Int    `json:"order"`
	BreakStartDatetime string      `json:"breakStartDatetime"`
	BreakEndDatetime   string      `json:"breakEndDatetime"`
	ActiveFlag         bool        `json:"activeFlag"`
}

func GetStaffsResponse(data *[]usecaseoutputs.Staff) ApiResponse {
	staffs := make([]Staff, 0)
	for _, t := range *data {
		reservation := Staff{
			StoreID:            t.StoreID,
			StaffID:            t.StaffID,
			Name:               t.Name,
			Order:              t.Order,
			BreakStartDatetime: utils.ConvertNTimeToString(t.BreakStartDatetime),
			BreakEndDatetime:   utils.ConvertNTimeToString(t.BreakStartDatetime),
			ActiveFlag:         t.ActiveFlag,
		}
		staffs = append(staffs, reservation)
	}

	return ApiResponse{
		Data:     staffs,
		Messages: []string{},
	}
}
