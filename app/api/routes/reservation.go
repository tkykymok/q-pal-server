package routes

import (
	"app/api/handlers"
	"app/pkg/reservation"
	"github.com/gofiber/fiber/v2"
)

func ReservationRouter(app fiber.Router, usecase reservation.Usecase) {
	app.Get("/reservations", handlers.GetReservationsByStoreId(usecase))
}
