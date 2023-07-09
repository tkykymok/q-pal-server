package message

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"sync"
)

var lock = &sync.Mutex{}

var fieldNamesMapInstance map[string]string

func initFieldNames() {
	if fieldNamesMapInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		file, _ := os.Open("/app/api/message/field-names.json")
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(file)

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&fieldNamesMapInstance); err != nil {
			log.Fatal(err)
		}
	}
}

func GetFieldName(field string) string {
	initFieldNames()
	return fieldNamesMapInstance[strings.ToLower(field)]
}
