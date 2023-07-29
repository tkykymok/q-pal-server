package message

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"
)

var messagesMapInstance map[string]string

func initMessages() {
	if messagesMapInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		file, _ := os.Open("/app/api/message/messages.json")
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(file)

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&messagesMapInstance); err != nil {
			log.Fatal(err)
		}
	}
}

func GetValidationMessage(errTag string, fieldName string) string {
	initMessages()
	message := messagesMapInstance[errTag]
	message = strings.Replace(message, "{fn}", GetFieldName(fieldName), 1)
	return message
}

func GetMessage(tag Tag, args ...string) string {
	initMessages()
	message := messagesMapInstance[string(tag)]
	for i, el := range args {
		message = strings.Replace(message, "{"+strconv.Itoa(i)+"}", el, 1)
	}
	return message
}
