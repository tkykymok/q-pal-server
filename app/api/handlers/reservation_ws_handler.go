package handlers

import (
	"app/api/presenter"
	"app/pkg/broadcast"
	"app/pkg/core/usecase"
	"context"
	"encoding/json"
	"log"
	"strconv"

	"app/api/requests"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func UpgradeReservationWsHandler(service usecase.ReservationUsecase) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return websocket.New(func(c *websocket.Conn) {
				storeIdStr := c.Query("storeId", "0")
				storeId, _ := strconv.Atoi(storeIdStr)

				// add the connection to the global clients map
				broadcast.ReservationClient.AddClient("reservation", storeId, c) // <- Use the broadcast package here
				defer func() {
					// remove the connection when done
					broadcast.ReservationClient.RemoveClient("reservation", storeId, c) // <- And here
					err := c.Close()
					if err != nil {
						return
					}
				}()
				handleReservationWsConnection(service)(c)
			})(c)
		}
		return c.Next()
	}
}

func handleReservationWsConnection(service usecase.ReservationUsecase) func(*websocket.Conn) {
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
