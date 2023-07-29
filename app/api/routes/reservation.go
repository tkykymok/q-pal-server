package routes

import (
	"app/api/handlers"
	"app/pkg/core/usecase"
	"github.com/gofiber/fiber/v2"
)

func ReservationRouter(app fiber.Router, usecase usecase.ReservationUsecase) {
	app.Get("/reservations", handlers.GetTodayReservations(usecase))
	app.Get("/lineEndWaitTime", handlers.GetLineEndWaitTime(usecase))
	app.Get("/individualWaitTime", handlers.GetIndividualWaitTime(usecase))
	app.Post("/reservation", handlers.CreateReservation(usecase))

	// WebSocket
	app.Get("/ws/reservation", handlers.UpgradeReservationWsHandler(usecase))
}
