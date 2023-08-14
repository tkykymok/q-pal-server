package routes

import (
	"app/api/handlers"
	"app/pkg/core/usecase"
	"github.com/gofiber/fiber/v2"
)

func StaffRouter(app fiber.Router, usecase usecase.StaffUsecase) {
	app.Get("/staffs", handlers.GetStaffs(usecase))
	app.Post("/create-active-staff", handlers.CreateActiveStaff(usecase))
}
