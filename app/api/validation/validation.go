package validation

import (
	"github.com/go-playground/validator/v10"
	"log"
	"sync"
)

var lock = &sync.Mutex{}

var validateInstance *validator.Validate

func initValidateInstance() {
	if validateInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		validateInstance = validator.New()

		// カスタムバリデーション登録
		registerCustomValidation("title-custom", customValidation)
	}
}

func registerCustomValidation(tag string, fn validator.Func) {
	err := validateInstance.RegisterValidation(tag, fn)
	if err != nil {
		log.Fatal(err)
	}
}

func ValidateStruct(request interface{}) error {
	initValidateInstance()
	err := validateInstance.Struct(request)
	if err != nil {
		return err
	}
	return nil
}
