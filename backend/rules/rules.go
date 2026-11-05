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
		if add > 400 || add < 1 {
			return "", fmt.Errorf("неправильная длина дня")
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
		needDay, err := parseMonthDays(rul[1])
		if err != nil {
			return "", err
		}
		for afterNow(check, now) {
			check = check.AddDate(0, 0, 1)
		}
		for !monthDayMatch(check, needDay) {
			check = check.AddDate(0, 0, 1)
		}

	default:
		return "", fmt.Errorf("недопустимый символ")
	}
	return check.Format("20060102"), nil
}

func afterNow(date, now time.Time) bool {
	return !date.After(now)
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

func parseMonthDays(s string) ([]int, error) {
	var res []int
	for _, p := range strings.Split(s, ",") {
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil, err
		}
		if n == 0 || n > 31 || n < -2 {
			return nil, fmt.Errorf("день месяца: от 1 до 31, -1 или -2")
		}
		res = append(res, n)
	}
	return res, nil
}

func monthDayMatch(date time.Time, days []int) bool {
	last := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
	for _, a := range days {
		if (a > 0 && date.Day() == a) || (a < 0 && date.Day() == last+a+1) {
			return true
		}
	}
	return false
}
