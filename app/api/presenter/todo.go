package presenter

import (
	"app/pkg/usecaseoutputs"
	"github.com/gofiber/fiber/v2"
	"github.com/volatiletech/null/v8"
)

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title,omitempty"`
	Completed bool   `json:"completed"`
	CreatedAt string `json:"createdAt,omitempty"`
}

type InputNotification struct {
	Type     string `json:"type"`
	IsTyping bool   `json:"is_typing"`
}

type TodoWithRelated struct {
	ID        int         `json:"id"`
	Title     string      `json:"title,omitempty"`
	Completed bool        `json:"completed"`
	Name      null.String `json:"name"`
	CreatedAt string      `json:"created_at,omitempty"`
}

func GetTodoByIdResponse(data *usecaseoutputs.Todo) *fiber.Map {
	todo := Todo{
		ID:        data.ID,
		Title:     data.Title,
		Completed: data.Completed,
		CreatedAt: data.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	return &fiber.Map{
		"data": todo,
	}
}

func GetAllTodosResponse(data *[]usecaseoutputs.Todo) *fiber.Map {
	todos := make([]Todo, 0)
	for _, t := range *data {
		todo := Todo{
			ID:        t.ID,
			Title:     t.Title,
			Completed: t.Completed,
			CreatedAt: t.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		todos = append(todos, todo)
	}

	return &fiber.Map{
		"todos": todos,
	}
}

func GetTodosWithRelatedResponse(data *[]usecaseoutputs.TodoWithRelated) *fiber.Map {
	todos := make([]TodoWithRelated, 0)
	for _, t := range *data {
		todo := TodoWithRelated{
			ID:        t.ID,
			Title:     t.Title,
			Completed: t.Completed,
			Name:      t.Name,
			CreatedAt: t.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		todos = append(todos, todo)
	}

	return &fiber.Map{
		"data": data,
	}
}
