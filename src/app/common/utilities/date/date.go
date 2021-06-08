package date

import (
	"time"
)

func StrToDate(t string) (time.Time, error) {
	layout := "2006-01-02"
	return time.Parse(layout, t)
}

func ValidDate(s string) bool {
	_, err := time.Parse("2006-01-02", s)
	if err != nil {
		return false
	}

	return true
}