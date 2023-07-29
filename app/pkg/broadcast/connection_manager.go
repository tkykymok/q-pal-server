package broadcast

import (
	"github.com/gofiber/websocket/v2"
	"log"
)

type ConnectionManager struct {
	clients map[string]map[int][]*websocket.Conn
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		clients: make(map[string]map[int][]*websocket.Conn),
	}
}

func (manager *ConnectionManager) AddClient(key string, subKey int, client *websocket.Conn) {
	if _, ok := manager.clients[key]; !ok {
		manager.clients[key] = make(map[int][]*websocket.Conn)
	}
	manager.clients[key][subKey] = append(manager.clients[key][subKey], client)
}

func (manager *ConnectionManager) RemoveClient(key string, subKey int, client *websocket.Conn) {
	if subMap, ok := manager.clients[key]; ok {
		if clients, ok := subMap[subKey]; ok {
			for i, c := range clients {
				if c == client {
					// Remove the client from the slice
					manager.clients[key][subKey] = append(clients[:i], clients[i+1:]...)
					break
				}
			}
			if len(manager.clients[key][subKey]) == 0 {
				delete(manager.clients[key], subKey)
			}
		}
		if len(manager.clients[key]) == 0 {
			delete(manager.clients, key)
		}
	}
}

func (manager *ConnectionManager) SendMessage(key string, subKey int, message interface{}) {
	// keyとsubKeyが指定されている場合は、それに紐づくクライアントにのみメッセージを送信
	if subMap, ok := manager.clients[key]; ok {
		if clients, ok := subMap[subKey]; ok {
			for _, client := range clients {
				err := client.WriteJSON(message)
				if err != nil {
					log.Printf("Error when broadcasting: %v", err)
					err := client.Close()
					if err != nil {
						return
					}
					manager.RemoveClient(key, subKey, client)
				}
			}
		}
	}
}
