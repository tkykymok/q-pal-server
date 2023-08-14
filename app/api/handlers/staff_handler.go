package handlers

import (
	"app/api/presenter"
	"app/pkg/core/usecase"
	"context"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// GetStaffs 店舗ごとのスタッフ一覧を取得する
func GetStaffs(usecase usecase.StaffUsecase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Create cancellable context.
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
