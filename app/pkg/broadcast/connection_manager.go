package broadcast

import "github.com/gofiber/websocket/v2"

type ConnectionManager struct {
	clients map[*websocket.Conn]bool
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		clients: make(map[*websocket.Conn]bool),
	}
}

func (manager *ConnectionManager) AddClient(client *websocket.Conn) {
	manager.clients[client] = true
}

func (manager *ConnectionManager) RemoveClient(client *websocket.Conn) {
	delete(manager.clients, client)
}
