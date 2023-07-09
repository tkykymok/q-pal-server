package requests

import "encoding/json"

type AddTodo struct {
	Title string `validate:"required,title-custom"`
}

type UpdateTodo struct {
	ID        int    `validate:"required"`
	Title     string `validate:"required,title-custom"`
	Completed bool
}

type GetTodosWithRelated struct {
	ID     int
	UserId int
}

type WSMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}
