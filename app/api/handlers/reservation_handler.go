package handlers

import (
	"app/api/presenter"
	"app/pkg/reservation"
	"context"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func GetReservationsByStoreId(usecase reservation.Usecase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Create cancellable context.
		customContext, cancel := context.WithCancel(context.Background())
		defer cancel()

		fetched, err := usecase.FetchReservationsByStoreId(customContext, 2)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.ErrorResponse(err))
		}
		return c.JSON(presenter.GetReservationsByStoreIdResponse(fetched))
	}
}
