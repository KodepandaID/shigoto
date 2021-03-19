package cronparser

import (
	"errors"
	"regexp"
)

// Regex to validate cron format
const (
	minute        = `^(\*|[1-5]?[0-9](-[1-5]?[0-9])?)(\/[1-5][0-9]*)?(,(\*|[1-5]?[0-9](-[1-5]?[0-9])?)(\/[1-9][0-9]*)?)*$`
	minuteValue   = `0?[0-9]|[1-5][0-9]`
	hour          = `^(\*|(1?[0-9]|2[0-3])(-(1?[0-9]|2[0-3]))?)(\/(1?[0-9]|2[0-3])(-(1?[0-9]|2[0-3]))?)?(,(\*|(1?[0-9]|2[0-3])(-(1?[0-9]|2[0-3]))?)(\/[1-9][0-9]*)?)*$`
	hourValue     = `0?[0-9]|1[0-9]|2[0-3]`
	dayMonth      = `^(\*|([1-9]|[1-2][0-9]?|3[0-1])(-([1-9]|[1-2][0-9]?|3[0-1]))?)(\/(1?[0-9]|2[0-9]|3[0-1]))?(,(\*|([1-9]|[1-2][0-9]?|3[0-1])(-([1-9]|[1-2][0-9]?|3[0-1]))?)?)*$`
	dayMonthValue = `0?[1-9]|[12][0-9]|3[01]`
	month         = `^(\*|([1-9]|1[0-2]?)(-([1-9]|1[0-2]?))?)(\/[1-9][0-9]*)?(,(\*|([1-9]|1[0-2]?)(-([1-9]|1[0-2]?))?)(\/[1-9][0-9]*)?)*$`
	monthValue    = `0?[1-9]|1[012]`
	weekday       = `^(?:MON|TUE|WED|THU|FRI|SAT|SUN)|(\*|[0-6](-[0-6])?)(\/[1-9][0-9]*)?(,(\*|[0-6](-[0-6])?)(\/[1-9][0-9]*)?)*$`
	weekdayValue  = `0?[0-7]|MON|TUE|WED|THU|FRI|SAT|SUN`
)

// validate cron format
func validate(expr []string) error {
	if len(expr) == 0 {
		return errors.New("Cron format cannot be empty string")
	}
	if len(expr) > 5 {
		return errors.New("Cron format is incorrect")
	}
	if validateMinute(expr[0]) == false {
		return errors.New("Cron format minute is incorrect")
	}
	if validateHour(expr[1]) == false {
		return errors.New("Cron format hour is incorrect")
	}
	if validateDayMonth(expr[2]) == false {
		return errors.New("Cron format day month is incorrect")
	}
	if validateMonth(expr[3]) == false {
		return errors.New("Cron format month is incorrect")
	}
	if validateWeekday(expr[4]) == false {
		return errors.New("Cron format weekday is incorrect")
	}

	return nil
}

func validateMinute(s string) bool {
	e := regexp.MustCompile(minute)
	match := e.FindAllString(s, -1)
	if len(match) == 0 {
		return false
	}

	return true
}

func validateHour(s string) bool {
	e := regexp.MustCompile(hour)
	match := e.FindAllString(s, -1)
	if len(match) == 0 {
		return false
	}

	return true
}

func validateDayMonth(s string) bool {
	e := regexp.MustCompile(dayMonth)
	match := e.FindAllString(s, -1)
	if len(match) == 0 {
		return false
	}

	return true
}

func validateMonth(s string) bool {
	e := regexp.MustCompile(month)
	match := e.FindAllString(s, -1)
	if len(match) == 0 {
		return false
	}

	return true
}

func validateWeekday(s string) bool {
	e := regexp.MustCompile(weekday)
	match := e.FindAllString(s, -1)
	if len(match) == 0 {
		return false
	}

	return true
}
