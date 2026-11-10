package rules

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// NextDate возвращает дату очередного повторения в формате 20060102.
// Поиск начинается после dstart; результат должен быть позже календарного дня now.
// Правило m поддерживает дни месяца (1–31, -1, -2) и необязательный список месяцев (1–12).
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	// Сохраняем календарный день now, задавая полночь UTC для сравнения с dstart.
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
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
		// Отсчитываем интервалы от исходной даты, чтобы сохранить периодичность.
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
		check = check.AddDate(0, 0, 1)
		for afterNow(check, now) {
			check = check.AddDate(0, 0, 1)
		}
		// После перехода за now ищем ближайший разрешённый день недели.
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
		// Формат: m <дни месяца> [<месяцы>], например "m 1,15" или "m 10,17 12,8,1".
		if len(rul) < 2 || len(rul) > 3 {
			return "", fmt.Errorf("неверная длина")
		}
		needDay, err := parseMonthDays(rul[1])
		if err != nil {
			return "", err
		}
		// Пустой список месяцев означает «любой месяц».
		var needMonth []int
		if len(rul) == 3 {
			needMonth, err = parseMonths(rul[2])
			if err != nil {
				return "", err
			}
		}
		// Ищем первую подходящую дату, которая позже и dstart, и now.
		check = check.AddDate(0, 0, 1)
		if afterNow(check, now) {
			check = now.AddDate(0, 0, 1)
		}
		// Перебор реальных дат пропускает отсутствующие в месяце числа, например 31 февраля.
		// Предел в 9 лет нужен, чтобы не зависнуть на правиле без решения («m 31 2»),
		// и при этом хватает для 29 февраля, которое может отсутствовать до 8 лет подряд.
		found := false
		for i := 0; i < 366*9; i++ {
			if monthDayMatch(check, needDay) && monthMatch(check, needMonth) {
				found = true
				break
			}
			check = check.AddDate(0, 0, 1)
		}
		if !found {
			return "", fmt.Errorf("по правилу m нет подходящей даты")
		}

	default:
		return "", fmt.Errorf("недопустимый символ")
	}
	return check.Format("20060102"), nil
}

// afterNow сообщает, что date раньше now или совпадает с ним.
func afterNow(date, now time.Time) bool {
	return !date.After(now)
}

// weekDay проверяет список дней недели и наличие в нём дня date.
// В правиле понедельнику соответствует 1, воскресенью — 7.
func weekDay(date time.Time, days []string) (bool, error) {
	// Проверяем весь список до поиска совпадения, чтобы не пропустить ошибочный элемент.
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

		// В time.Weekday воскресенье обозначается нулём, а в правиле — семёркой.
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

// parseMonthDays разбирает дни месяца: 1–31, -1 (последний), -2 (предпоследний).
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

// monthDayMatch проверяет совпадение даты с одним из заданных дней месяца.
func monthDayMatch(date time.Time, days []int) bool {
	// Нулевой день следующего месяца определяет последний день текущего.
	// Для -1 и -2 выражение last+a+1 выбирает последний и предпоследний день.
	last := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
	for _, a := range days {
		if (a > 0 && date.Day() == a) || (a < 0 && date.Day() == last+a+1) {
			return true
		}
	}
	return false
}

// parseMonths разбирает список месяцев: числа от 1 до 12.
func parseMonths(s string) ([]int, error) {
	var res []int
	for _, p := range strings.Split(s, ",") {
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil, err
		}
		if n < 1 || n > 12 {
			return nil, fmt.Errorf("месяц должен быть от 1 до 12")
		}
		res = append(res, n)
	}
	return res, nil
}

// monthMatch проверяет, что месяц date входит в список; пустой список подходит любому месяцу.
func monthMatch(date time.Time, months []int) bool {
	if len(months) == 0 {
		return true
	}
	for _, m := range months {
		if int(date.Month()) == m {
			return true
		}
	}
	return false
}
