package handlers

import (
	"app/api/presenter"
	"app/pkg/core/usecase"
	"context"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

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
