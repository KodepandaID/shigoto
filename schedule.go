package shigoto

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

// CronFormat to set a job with cron format
func (j *Jobs) CronFormat(cron string) *Jobs {
	cronSlice := strings.Split(cron, " ")
	j.Cron = cronSlice

	return j
}

// At to run a job at a time
func (j *Jobs) At(time string) {
	parts := strings.Split(time, ":")
	if len(parts) < 2 {
		log.Fatal("The clock format is wrong")
	}

	hour, _ := strconv.Atoi(parts[0])
	minute, _ := strconv.Atoi(parts[1])

	j.Cron[0] = fmt.Sprintf("%d", minute)
	j.Cron[1] = fmt.Sprintf("%d", hour)
}

// EveryMinute to run a job every minutes
func (j *Jobs) EveryMinute() *Jobs {
	j.Cron = []string{"*", "*", "*", "*", "*"}
	return j
}

// EveryFiveMinutes to run a jobs every 5 minutes
func (j *Jobs) EveryFiveMinutes() *Jobs {
	j.Cron = []string{"*/5", "*", "*", "*", "*"}
	return j
}

// EveryTenMinutes to run a job every 10 minutes
func (j *Jobs) EveryTenMinutes() *Jobs {
	j.Cron = []string{"*/10", "*", "*", "*", "*"}
	return j
}

// EveryFifteenMinutes to run a job every 15 minutes
func (j *Jobs) EveryFifteenMinutes() *Jobs {
	j.Cron = []string{"*/15", "*", "*", "*", "*"}
	return j
}

// EveryThirtyMinutes to run a job every 30 minutes
func (j *Jobs) EveryThirtyMinutes() *Jobs {
	j.Cron = []string{"*/30", "*", "*", "*", "*"}
	return j
}

// Hourly to run a job every hours
func (j *Jobs) Hourly() *Jobs {
	j.Cron = []string{"0", "*/1", "*", "*", "*"}
	return j
}

// EveryThreeHours to run a job every 3 hours
func (j *Jobs) EveryThreeHours() *Jobs {
	j.Cron = []string{"0", "*/3", "*", "*", "*"}
	return j
}

// EverySixHours to run a job every 6 hours
func (j *Jobs) EverySixHours() *Jobs {
	j.Cron = []string{"0", "*/6", "*", "*", "*"}
	return j
}

// EveryTwelveHours to run a job every 12 hours
func (j *Jobs) EveryTwelveHours() *Jobs {
	j.Cron = []string{"0", "*/12", "*", "*", "*"}
	return j
}

// Daily to run a job every day at midnight
func (j *Jobs) Daily() *Jobs {
	j.Cron = []string{"0", "0", "*/1", "*", "*"}
	return j
}

// DailyAt to run a job every day at a specific time
func (j *Jobs) DailyAt(time string) *Jobs {
	parts := strings.Split(time, ":")
	if len(parts) < 2 {
		log.Fatal("The clock format is wrong")
	}

	hour, _ := strconv.Atoi(parts[0])
	minute, _ := strconv.Atoi(parts[1])

	j.Cron = []string{fmt.Sprint(minute), fmt.Sprint(hour), "*/1", "*", "*"}
	return j
}

// Weekly to run a job every week
func (j *Jobs) Weekly() *Jobs {
	j.Cron = []string{"0", "0", "*/7", "*", "*"}
	return j
}

// WeeklyOn to run a job every week at a specific time
func (j *Jobs) WeeklyOn(time string) *Jobs {
	parts := strings.Split(time, ":")
	if len(parts) < 2 {
		log.Fatal("The clock format is wrong")
	}

	hour, _ := strconv.Atoi(parts[0])
	minute, _ := strconv.Atoi(parts[1])

	j.Cron = []string{fmt.Sprint(minute), fmt.Sprint(hour), "*/7", "*", "*"}
	return j
}

// Monthly to run a job every month
func (j *Jobs) Monthly() *Jobs {
	j.Cron = []string{"0", "0", "1", "*/1", "*"}
	return j
}

// MonthlyOn to run a job every month at a specific time
func (j *Jobs) MonthlyOn(time string) *Jobs {
	parts := strings.Split(time, ":")
	if len(parts) < 2 {
		log.Fatal("The clock format is wrong")
	}

	hour, _ := strconv.Atoi(parts[0])
	minute, _ := strconv.Atoi(parts[1])

	j.Cron = []string{fmt.Sprint(minute), fmt.Sprint(hour), "1", "*/1", "*"}
	return j
}

// Quarterly to run a job every 6 month
func (j *Jobs) Quarterly() *Jobs {
	j.Cron = []string{"0", "0", "1", "*/6", "*"}
	return j
}

// Yearly to run a job every 1 year
func (j *Jobs) Yearly() *Jobs {
	j.Cron = []string{"0", "0", "1", "*/12", "*"}
	return j
}
