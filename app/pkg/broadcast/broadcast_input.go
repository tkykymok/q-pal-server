package broadcast

import (
	"app/api/presenter"
)

var TodoInputClient = NewInputConnectionManager()

type InputConnectionManager struct {
	*ConnectionManager
}

func NewInputConnectionManager() *InputConnectionManager {
	return &InputConnectionManager{
		ConnectionManager: NewConnectionManager(),
	}
}

func (manager *InputConnectionManager) SendInputNotification(notification presenter.InputNotification) {
	manager.SendMessage("todo_input", 1, notification)
}
