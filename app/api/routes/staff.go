package routes

import (
	"app/api/handlers"
	"app/pkg/core/usecase"
	"github.com/gofiber/fiber/v2"
)

func StaffRouter(app fiber.Router, usecase usecase.StaffUsecase) {
	app.Get("/staffs", handlers.GetStaffs(usecase))
	app.Post("/create-active-staff", handlers.CreateActiveStaff(usecase))
	app.Put("/update-active-staffs", handlers.UpdateActiveStaffs(usecase))
	app.Delete("/remove-active-staff/:staffId", handlers.RemoveActiveStaff(usecase))
}
