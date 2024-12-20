package task

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const layoutDate = "20060102"

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("repeat is empty")
	}

	validDate, err := time.Parse(layoutDate, date)
	if err != nil {
		return "", fmt.Errorf("incorrect date %v", err)
	}

	rule := string(repeat[0])
	rightLen := len(repeat) > 2
	var result string

	switch {
	//задача переносится на указанное число дней
	case rule == "d" && rightLen:
		result, err = everyDay(now, validDate, repeat[2:])

	// задача назначается в указанные дни недели, где 1 — понедельник, 7 — воскресенье
	case rule == "w" && rightLen:
		result, err = everyWeek(validDate, now, repeat[2:])

	// задача назначается в указанные дни месяца (1-31)
	case rule == "m" && rightLen:
		result, err = everyMonth(validDate, now, repeat[2:])

	// задача выполняется ежегодно
	case rule == "y":
		result, err = everyYear(now, validDate)
	default:
		return "", fmt.Errorf("incorrect repetition rule %v", err)
	}

	return result, err
}

func everyDay(now, date time.Time, days string) (string, error) {
	d, err := strconv.Atoi(days)
	if err != nil || d > 400 || d < 0 {
		return "", fmt.Errorf(`incorrect repetition rule in "d"`)
	}

	resultDate := date.AddDate(0, 0, d)
	for resultDate.Before(now) {
		resultDate = resultDate.AddDate(0, 0, d)
	}

	return resultDate.Format(layoutDate), nil
}

func everyWeek(date, now time.Time, repeat string) (string, error) {
	result := ""

	if date.Before(now) {
		date = now
	}

	days := strings.Split(repeat, ",")
	week := make(map[int]string)

	for i := 1; i <= 7; i++ {
		date = date.AddDate(0, 0, 1)
		weekDay := int(date.Weekday())

		if weekDay == 0 {
			weekDay = 7
		}

		week[weekDay] = date.Format(layoutDate)

		for _, day := range days {
			d, err := strconv.Atoi(day)
			if err != nil || d > 7 || d < 0 {
				return "", fmt.Errorf(`incorrect repetition rule in "w" %v`, err)
			}

			if d == weekDay {
				result = week[d]
				return result, nil
			}
		}
	}

	return result, nil
}

func everyMonth(date, now time.Time, repeat string) (string, error) {
	result := ""

	if date.Before(now) {
		date = now
	}

	// получаем количество аргументов repeat
	args := strings.Split(repeat, " ")

	// первый аргумент - по каким дням
	days := strings.Split(args[0], ",")

	if len(args) == 1 {
		needDate, err := onlyDays(date, days)
		if err != nil {
			return "", err
		}

		result = needDate.Format(layoutDate)
	}

	// второй аргумент - по каким месяцам
	if len(args) > 1 {
		months := strings.Split(args[1], ",")

		needDate, err := monthAndDays(date, months, days)
		if err != nil {
			return "", err
		}

		result = needDate.Format(layoutDate)
	}

	return result, nil
}

func onlyDays(date time.Time, days []string) (time.Time, error) {
	month := getNextMonth(date)

	resultSlice := make([]time.Time, 0)

	for _, v := range days {
		targetDay, err := strconv.Atoi(v)
		if err != nil || targetDay > 31 || targetDay == 0 || targetDay < -2 {
			return date, fmt.Errorf(`incorrect repetition rule in "m"`)
		}

		date = date.AddDate(0, 0, 1)

		for _, day := range month {
			ifTargetDayNegative(&targetDay, day)

			if targetDay == int(day.Day()) {
				resultSlice = append(resultSlice, day)
			}
		}
	}

	return resultDate(date, &resultSlice), nil
}

// функция для получения следующих двух месяцев
func getNextMonth(date time.Time) []time.Time {
	month := make([]time.Time, 0)

	day := date.AddDate(0, 0, 1)

	for j := 0; j < 62; j++ {
		month = append(month, day)

		day = day.AddDate(0, 0, 1)
	}

	return month
}

func monthAndDays(date time.Time, month, days []string) (time.Time, error) {
	year := getNextYear(date)

	resultSlice := make([]time.Time, 0)

	for _, m := range month {
		targetMonth, err := strconv.Atoi(m)
		if err != nil || targetMonth > 12 || targetMonth <= 0 {
			return date, fmt.Errorf(`incorrect repetition rule in "m"`)
		}

		date = date.AddDate(0, 0, 1)

		for _, day := range year[targetMonth] {

			for _, d := range days {
				targetDay, err := strconv.Atoi(d)
				if err != nil || targetDay > 31 || targetDay < -2 {
					return date, fmt.Errorf(`incorrect repetition rule in "m"`)
				}

				ifTargetDayNegative(&targetDay, day)

				if targetDay == int(day.Day()) {
					resultSlice = append(resultSlice, day)
				}
			}

		}
	}

	return resultDate(date, &resultSlice), nil
}

func getNextYear(date time.Time) map[int][]time.Time {
	year := make(map[int][]time.Time)
	var month time.Month
	day := date.AddDate(0, 0, 1)

	for i := 0; i < 12; i++ {
		month = day.Month()

		for j := 0; j < 31; j++ {
			year[int(month)] = append(year[int(month)], day)

			if day == day.AddDate(0, 1, -day.Day()) {
				break
			}

			day = day.AddDate(0, 0, 1)

		}
		day = day.AddDate(0, 0, 1)
	}

	return year
}

func ifTargetDayNegative(targetDay *int, day time.Time) {
	switch *targetDay {
	case -1:
		day = day.AddDate(0, 1, -day.Day())
		*targetDay = int(day.Day())
	case -2:
		day = day.AddDate(0, 1, -day.Day()-1)
		*targetDay = int(day.Day())
	}
}

func resultDate(date time.Time, resultSlice *[]time.Time) time.Time {
	for {
		date = date.AddDate(0, 0, 1)

		for _, v := range *resultSlice {
			if date.Truncate(24 * time.Hour).Equal(v.Truncate(24 * time.Hour)) {
				return v
			}
		}
	}
}

func everyYear(now, date time.Time) (string, error) {
	if date.Before(now) {
		for date.Before(now) {
			date = date.AddDate(1, 0, 0)
		}
	} else {
		date = date.AddDate(1, 0, 0)
	}

	return date.Format(layoutDate), nil
}
