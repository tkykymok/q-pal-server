package presenter

import "github.com/gofiber/fiber/v2"

func GetSuccessResponse(messages ...string) *fiber.Map {
	return &fiber.Map{
		"messages": messages,
	}
}
