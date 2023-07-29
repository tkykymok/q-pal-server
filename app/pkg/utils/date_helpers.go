package utils

import (
	"os"
	"time"
)

func GetStartOfDay() time.Time {
	location, _ := time.LoadLocation(os.Getenv("LOC"))
	currentTime := time.Now().In(location)
	startOfDay := time.Date(
		currentTime.Year(),
		currentTime.Month(),
		currentTime.Day(),
		0, 0, 0, 0,
		location,
	)
	return startOfDay
}

func GetEndOfDay() time.Time {
	location, _ := time.LoadLocation(os.Getenv("LOC"))
	currentTime := time.Now().In(location)
	endOfDay := time.Date(
		currentTime.Year(),
		currentTime.Month(),
		currentTime.Day(),
		23, 59, 59, 0,
		location,
	)
	return endOfDay
}
