package routes

import (
	"app/api/handlers"
	"app/pkg/core/usecase"
	"github.com/gofiber/fiber/v2"
)

func ReservationRouter(app fiber.Router, usecase usecase.ReservationUsecase) {
	app.Get("/reservations/today", handlers.GetTodayReservations(usecase))
	app.Get("/reservations/line-end-wait-time", handlers.GetLineEndWaitTime(usecase))
	app.Get("/reservations/my-wait-time", handlers.GetMyWaitTime(usecase))
	app.Post("/create-reservation", handlers.CreateReservation(usecase))
	app.Put("/update-reservation/status", handlers.UpdateReservationStatus(usecase))

	// WebSocket
	app.Get("/ws/reservations", handlers.UpgradeReservationWsHandler(usecase))
}
