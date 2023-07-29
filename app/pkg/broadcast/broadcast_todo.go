package broadcast

import (
	"app/api/presenter"
)

var TodoClient = NewTodoConnectionManager()

type TodoConnectionManager struct {
	*ConnectionManager
}

func NewTodoConnectionManager() *TodoConnectionManager {
	return &TodoConnectionManager{
		ConnectionManager: NewConnectionManager(),
	}
}

func (manager *TodoConnectionManager) SendNewTodo(todo presenter.Todo) {
	manager.SendMessage("todo", 1, todo)
}
