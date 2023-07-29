package handlers

import (
	"app/api/requests"
	"app/pkg/broadcast"
	"app/pkg/core/usecase"
	"context"
	"encoding/json"
	"log"
	"strconv"

	"app/api/presenter"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func UpgradeTodoInputWsHandler(service usecase.TodoUsecase) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return websocket.New(func(c *websocket.Conn) {
				userIdStr := c.Query("userId", "0")
				userId, _ := strconv.Atoi(userIdStr)

				// add the connection to the global clients map
				broadcast.TodoInputClient.AddClient("todo_input", userId, c) // <- Use the broadcast package here
				defer func() {
					// remove the connection when done
					broadcast.TodoInputClient.RemoveClient("todo_input", userId, c) // <- And here
					err := c.Close()
					if err != nil {
						return
					}
				}()
				handleTodoInputWsConnection(service)(c)
			})(c)
		}
		return c.Next()
	}
}

func handleTodoInputWsConnection(service usecase.TodoUsecase) func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		// Create cancellable context.
		_, cancel := context.WithCancel(context.Background())
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
				"on input": func() error {
					var notification presenter.InputNotification
					err := json.Unmarshal(wsMsg.Data, &notification)
					if err != nil {
						return err
					}
					broadcast.TodoInputClient.SendInputNotification(notification)
					return nil
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
