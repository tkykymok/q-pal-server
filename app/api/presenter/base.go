package presenter

import "github.com/gofiber/fiber/v2"

type BaseResponse struct {
	Messages []string `json:"messages" `
}

func GetSuccessResponse(messages ...string) *fiber.Map {
	data := BaseResponse{
		Messages: messages,
	}

	return &fiber.Map{
		"data": data,
	}
}
