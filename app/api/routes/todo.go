package routes

import (
	"app/api/handlers"
	"app/pkg/core/usecase"
	"github.com/gofiber/fiber/v2"
)

func TodoRouter(app fiber.Router, usecase usecase.TodoUsecase) {
	app.Get("/todos", handlers.GetAllTodos(usecase))
	app.Get("/todosWithRelated", handlers.GetTodosWithRelated(usecase))
	app.Get("/todo/:id", handlers.GetTodoById(usecase))
	app.Post("/todo", handlers.AddTodo(usecase))
	app.Put("/todo", handlers.UpdateTodo(usecase))

	// Add this line
	app.Get("/ws/todo", handlers.UpgradeTodoWsHandler(usecase))
	app.Get("/ws/todo/input", handlers.UpgradeTodoInputWsHandler(usecase))
}
