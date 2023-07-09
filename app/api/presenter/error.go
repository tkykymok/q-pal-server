package presenter

import (
	"app/api/message"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ValidationError struct {
	FailedField string `json:"failedField"`
	Tag         string `json:"tag"`
	Value       string `json:"value"`
	Message     string `json:"message"`
}

func ErrorResponse(err error) *fiber.Map {
	return &fiber.Map{
		"status": false,
		"data":   "",
		"error":  err.Error(),
	}
}

func ValidationErrorResponse(err error) *fiber.Map {
	var errors []*ValidationError
	for _, err := range err.(validator.ValidationErrors) {
		var temp ValidationError
		temp.FailedField = err.StructField()
		temp.Tag = err.Tag()
		temp.Value = err.Param()
		temp.Message = message.GetValidationMessage(err.Tag(), err.StructField())
		errors = append(errors, &temp)
	}

	return &fiber.Map{
		"errors": errors,
	}
}
