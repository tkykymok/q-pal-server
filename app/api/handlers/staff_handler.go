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

func CreateActiveStaff(usecase usecase.StaffUsecase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		customContext, cancel := context.WithCancel(context.Background())
		defer cancel()

		var request requests.AddActiveStaff
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
		return c.JSON(presenter.GetSuccessResponse(message.GetMessage(message.SUCCESS, "アクティブスタッフ登録")))
	}
}

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
		return c.JSON(presenter.GetSuccessResponse(message.GetMessage(message.SUCCESS, "アクティブスタッフ削除")))
	}
}
