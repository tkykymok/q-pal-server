package validation

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

func customValidation(fl validator.FieldLevel) bool {
	str, _ := fl.Field().Interface().(string)
	reg := regexp.MustCompile("[?!@＠]")
	return !reg.MatchString(str)
}
