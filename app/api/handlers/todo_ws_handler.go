package handlers

import (
	"app/pkg/broadcast"
	"app/pkg/core/usecase"
	"context"
	"encoding/json"
	"log"

	"app/api/presenter"
	"app/api/requests"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func UpgradeTodoWsHandler(service usecase.TodoUsecase) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return websocket.New(func(c *websocket.Conn) {
				// add the connection to the global clients map
				broadcast.TodoClient.AddClient(c) // <- Use the broadcast package here
				defer func() {
					// remove the connection when done
					broadcast.TodoClient.RemoveClient(c) // <- And here
					err := c.Close()
					if err != nil {
						return
					}
				}()
				handleTodoWsConnection(service)(c)
			})(c)
		}
		return c.Next()
	}
}

func handleTodoWsConnection(service usecase.TodoUsecase) func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		// Create cancellable context.
		customContext, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Infinite loop to handle multiple messages
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}

			var wsMsg requests.WSMessage
			err = json.Unmarshal(msg, &wsMsg)
			if err != nil {
				log.Println("message parse error:", err)
				break
			}

			// Map of message types to handlers
			handlers := map[string]func() error{
				"get all todos": func() error {
					fetched, err := service.FetchAllTodos(customContext)
					if err != nil {
						return err
					}
					response := presenter.GetAllTodosResponse(fetched)
					return c.WriteJSON(response)
				},
				"get todos with related": func() error {
					var request requests.GetTodosWithRelated
					err := json.Unmarshal(wsMsg.Data, &request)
					if err != nil {
						return err
					}
					fetched, err := service.FetchTodosWithRelated(customContext, &request)
					if err != nil {
						return err
					}
					response := presenter.GetTodosWithRelatedResponse(fetched)
					return c.WriteJSON(response)
				},
			}

			handler, ok := handlers[wsMsg.Type]
			if !ok {
				log.Println("unknown message type:", wsMsg.Type)
				continue
			}

			err = handler()
			if err != nil {
				log.Println("handler error:", err)
			}
		}
	}
}
