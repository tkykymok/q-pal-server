package handlers

import (
	"app/api/message"
	"app/api/presenter"
	"app/api/requests"
	"app/pkg/core/usecase"
	"app/pkg/usecaseinputs"
	"context"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// GetStaffs 店舗ごとのスタッフ一覧を取得する
func GetStaffs(usecase usecase.StaffUsecase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		customContext, cancel := context.WithCancel(context.Background())
		defer cancel()

		fetched, err := usecase.FetchStaffs(customContext, 2)

		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.ErrorResponse(err))
		}
		return c.JSON(presenter.GetStaffsResponse(fetched))
	}
}

// CreateActiveStaff 対応スタッフを登録する
func CreateActiveStaff(usecase usecase.StaffUsecase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		customContext, cancel := context.WithCancel(context.Background())
		defer cancel()

		var request requests.CreateActiveStaff
		err := c.BodyParser(&request)
		if err != nil {
			return err
		}

		input := usecaseinputs.CreateActiveStaffInput{
			StoreID: 2,
			StaffID: request.StaffId,
		}

		err = usecase.CreateActiveStaff(customContext, input)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.ErrorResponse(err))
		}
		return c.JSON(presenter.GetSuccessResponse(message.GetMessage(message.SUCCESS, "対応スタッフの登録")))
	}
}

// UpdateActiveStaffs 対応スタッフを更新する
func UpdateActiveStaffs(usecase usecase.StaffUsecase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		customContext, cancel := context.WithCancel(context.Background())
		defer cancel()

		var request requests.UpdateActiveStaff
		err := c.BodyParser(&request)
		if err != nil {
			return err
		}

		data := make([]usecaseinputs.UpdateActiveStaffData, 0)
		for _, t := range request.Data {
			temp := usecaseinputs.UpdateActiveStaffData{
				StaffID: t.StaffId,
				Order:   t.Order,
			}

			data = append(data, temp)
		}

		input := usecaseinputs.UpdateActiveStaffInput{
			StoreId: 2,
			Data:    data,
		}

		err = usecase.UpdateActiveStaff(customContext, input)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.ErrorResponse(err))
		}
		return c.JSON(presenter.GetSuccessResponse(message.GetMessage(message.SUCCESS, "対応スタッフの更新")))
	}
}

// RemoveActiveStaff 対応スタッフを削除する
func RemoveActiveStaff(usecase usecase.StaffUsecase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		customContext, cancel := context.WithCancel(context.Background())
		defer cancel()

		staffId, _ := c.ParamsInt("staffId", 0)

		input := usecaseinputs.RemoveActiveStaffInput{
			StoreID: 2,
			StaffID: staffId,
		}

		err := usecase.RemoveActiveStaff(customContext, input)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.ErrorResponse(err))
		}
		return c.JSON(presenter.GetSuccessResponse(message.GetMessage(message.SUCCESS, "対応スタッフの削除")))
	}
}
