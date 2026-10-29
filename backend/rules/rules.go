package rules

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("не может быть пустой строки")
	}
	rul := strings.Split(repeat, " ")
	check, err := time.Parse("20060102", dstart)

	var nexttime string
	switch rul[0] {
	case "d":
		add, err := strconv.Atoi(rul[1])
		if err != nil {
			return "", err
		}
		check.AddDate(0, 0, add)
		for !afterNow(check, now) {
			check.AddDate(0, 0, 1)
		}

	case "y":
		check.AddDate(1, 0, 0)
		for !afterNow(check, now) {
			check.AddDate(1, 0, 0)
		}
	case "w":

	case "m":

	default:
		return "удали задачу", err
	}

	if err != nil {
		return "", err
	}
	return check.Format("20060102"), nil
}

func afterNow(date, now time.Time) bool {
	return time.Now().After(date)
}
