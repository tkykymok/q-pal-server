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

// GetTodayReservations 店舗ごとの今日の予約一覧を取得する
func GetTodayReservations(usecase usecase.ReservationUsecase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Create cancellable context.
		customContext, cancel := context.WithCancel(context.Background())
		defer cancel()

		fetched, err := usecase.FetchAllReservations(customContext, 2)

		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.ErrorResponse(err))
		}
		return c.JSON(presenter.GetReservationsResponse(fetched))
	}
}

// GetLineEndWaitTime 店舗ごとの最後尾の 予約番号 & 待ち時間 を取得する
func GetLineEndWaitTime(usecase usecase.ReservationUsecase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Create cancellable context.
		customContext, cancel := context.WithCancel(context.Background())
		defer cancel()

		fetched, err := usecase.FetchLineEndWaitTime(customContext, 2)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.ErrorResponse(err))
		}
		return c.JSON(presenter.GetWaitTimeResponse(fetched))
	}
}

// GetIndividualWaitTime 顧客ごとの現在の待ち時間を取得する
func GetIndividualWaitTime(usecase usecase.ReservationUsecase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Create cancellable context.
		customContext, cancel := context.WithCancel(context.Background())
		defer cancel()

		fetched, err := usecase.FetchIndividualWaitTime(customContext, 2, 17)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.ErrorResponse(err))
		}
		return c.JSON(presenter.GetWaitTimeResponse(fetched))
	}
}

func CreateReservation(usecase usecase.ReservationUsecase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		customContext, cancel := context.WithCancel(context.Background())
		defer cancel()

		var request requests.CreateReservations
		err := c.BodyParser(&request)
		if err != nil {
			return err
		}

		input := usecaseinputs.CreateReservationInput{
			CustomerID: 0,
			StoreID:    request.StoreId,
			MenuID:     1, // TODO
		}

		err = usecase.CreateReservation(customContext, &input)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.ErrorResponse(err))
		}

		return c.JSON(presenter.GetSuccessResponse(message.GetMessage(message.SUCCESS, "予約")))
	}
}
