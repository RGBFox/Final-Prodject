package rules

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("необходимо удалить задачу")
	}
	rul := strings.Split(repeat, " ")
	check, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", err
	}

	switch rul[0] {
	case "d":
		if len(rul) != 2 {
			return "", fmt.Errorf("неверная длина")
		}
		add, err := strconv.Atoi(rul[1])
		if err != nil {
			return "", err
		}
		if add > 400 {
			return "", fmt.Errorf("слишком большая длина дня")
		}
		check = check.AddDate(0, 0, add)
		for afterNow(check, now) {
			check = check.AddDate(0, 0, add)
		}

	case "y":
		check = check.AddDate(1, 0, 0)
		for afterNow(check, now) {
			check = check.AddDate(1, 0, 0)
		}

	case "w":
		if len(rul) != 2 {
			return "", fmt.Errorf("неверная длина")
		}
		needDay := strings.Split(rul[1], ",")
		for afterNow(check, now) {
			check = check.AddDate(0, 0, 1)
		}
		for ok, err := weekDay(check, needDay); !ok; ok, _ = weekDay(check, needDay) {
			if err != nil {
				return "", err
			}
			ok, err = weekDay(check, needDay)
			if err != nil {
				return "", err
			}
			check = check.AddDate(0, 0, 1)
		}
	case "m":
		if len(rul) != 2 {
			return "", fmt.Errorf("неверная длина")
		}
		needDay := strings.Split(rul[1], ",")
		for afterNow(check, now) {
			check = check.AddDate(0, 0, 1)
		}
		for ok, err := mounthDay(check, needDay); !ok; ok, _ = mounthDay(check, needDay) {
			if err != nil {
				return "", err
			}
			ok, err = mounthDay(check, needDay)
			if err != nil {
				return "", err
			}
			check = check.AddDate(0, 0, 1)
		}

	default:
		return "недопустимый символ", nil
	}

	if err != nil {
		return "", err
	}
	return check.Format("20060102"), nil
}

func afterNow(date, now time.Time) bool {
	return now.After(date)
}

func weekDay(date time.Time, days []string) (bool, error) {

	for _, b := range days {
		a, err := strconv.Atoi(b)
		if err != nil {
			return false, err
		}

		if a < 1 || a > 7 {
			return false, fmt.Errorf("день недели должен быть от 1 до 7")
		}
	}

	for _, b := range days {
		a, err := strconv.Atoi(b)
		if err != nil {
			return false, err
		}

		wd := int(date.Weekday())
		if wd == 0 {
			wd = 7
		}
		if wd == a {
			return true, nil
		}
	}
	return false, nil
}

func mounthDay(date time.Time, days []string) (bool, error) {

	for _, b := range days {
		a, err := strconv.Atoi(b)
		if err != nil {
			return false, err
		}
		if a > 31 || a < -2 || a == 0 {
			return false, fmt.Errorf("в месяце не может быть больше 31 дня или меньше -2")
		}
	}

	for _, b := range days {
		a, err := strconv.Atoi(b)
		if err != nil {
			return false, err
		}
		if int(date.Day()) == a {
			return true, nil
		}
	}
	return false, nil
}
