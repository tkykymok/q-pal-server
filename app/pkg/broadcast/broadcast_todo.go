package broadcast

import (
	"app/api/presenter"
	"log"
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
	for client := range manager.clients {
		err := client.WriteJSON(todo)
		if err != nil {
			log.Printf("Error when broadcasting: %v", err)
			err := client.Close()
			if err != nil {
				return
			}
			manager.RemoveClient(client)
		}
	}
}
