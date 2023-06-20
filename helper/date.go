package helper

import (
	"time"
)

func StringToDate(dateString string) (time.Time, error) {
	date, err := time.Parse("2006-01-02", dateString)

	return date, err
}

func DatePassword(dateString string) (string, error) {
	date, err := time.Parse("2006-01-02", dateString)

	return date.Format("20060102"), err
}
