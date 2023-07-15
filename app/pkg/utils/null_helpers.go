package utils

import (
	"github.com/volatiletech/null/v8"
	"time"
)

func ConvertNTimeToString(t null.Time) string {
	if t.Time.IsZero() {
		return ""
	}
	return t.Time.Format("2006-01-02 15:04:05")
}

func ConvertTimeToString(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
