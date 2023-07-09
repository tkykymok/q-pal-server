package broadcast

import (
	"app/api/presenter"
	"log"
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
	for client := range manager.clients {
		err := client.WriteJSON(notification)
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
