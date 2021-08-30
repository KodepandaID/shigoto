package test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/KodepandaID/shigoto"
	"github.com/joho/godotenv"
)

var client *shigoto.Config

func init() {
	wd, _ := os.Getwd()
	p := filepath.Dir(wd)

	e := godotenv.Load(p + "/.env")
	if e != nil {
		os.Setenv("MONGO_URI", "mongodb://localhost:27017") // for github actions
	}

	client, _ = shigoto.New(&shigoto.Config{
		DB:     os.Getenv("MONGO_URI"),
		DBName: "jobs-scheduler",
	})
}

func TestCronFormat(t *testing.T) {
	client.Register("schedule-test", scheduleTest)
	jobs := client.Command("cron-format", "schedule-test")
	jobs.CronFormat("* * * * *")

	if !reflect.DeepEqual(jobs.Cron, []string{"*", "*", "*", "*", "*"}) {
		t.Fail()
	}
}

func TestAt(t *testing.T) {
	client.Register("schedule-test", scheduleTest)
	jobs := client.Command("cron-format", "schedule-test")
	jobs.At("01:00")

	if !reflect.DeepEqual(jobs.Cron, []string{"0", "1", "*", "*", "*"}) {
		t.Fail()
	}
}

func TestEveryMinute(t *testing.T) {
	client.Register("schedule-test", scheduleTest)
	jobs := client.Command("cron-format", "schedule-test")
	jobs.EveryMinute()

	if !reflect.DeepEqual(jobs.Cron, []string{"*", "*", "*", "*", "*"}) {
		t.Fail()
	}
}

func TestEveryFiveMinute(t *testing.T) {
	client.Register("schedule-test", scheduleTest)
	jobs := client.Command("cron-format", "schedule-test")
	jobs.EveryFiveMinutes()

	if !reflect.DeepEqual(jobs.Cron, []string{"*/5", "*", "*", "*", "*"}) {
		t.Fail()
	}
}

func TestEveryTenMinute(t *testing.T) {
	client.Register("schedule-test", scheduleTest)
	jobs := client.Command("cron-format", "schedule-test")
	jobs.EveryTenMinutes()

	if !reflect.DeepEqual(jobs.Cron, []string{"*/10", "*", "*", "*", "*"}) {
		t.Fail()
	}
}

func TestEveryFifteenMinute(t *testing.T) {
	client.Register("schedule-test", scheduleTest)
	jobs := client.Command("cron-format", "schedule-test")
	jobs.EveryFifteenMinutes()

	if !reflect.DeepEqual(jobs.Cron, []string{"*/15", "*", "*", "*", "*"}) {
		t.Fail()
	}
}

func TestEveryThirtyMinute(t *testing.T) {
	client.Register("schedule-test", scheduleTest)
	jobs := client.Command("cron-format", "schedule-test")
	jobs.EveryThirtyMinutes()

	if !reflect.DeepEqual(jobs.Cron, []string{"*/30", "*", "*", "*", "*"}) {
		t.Fail()
	}
}

func TestHourly(t *testing.T) {
	client.Register("schedule-test", scheduleTest)
	jobs := client.Command("cron-format", "schedule-test")
	jobs.Hourly()

	if !reflect.DeepEqual(jobs.Cron, []string{"0", "*/1", "*", "*", "*"}) {
		t.Fail()
	}
}

func TestEveryThreeHours(t *testing.T) {
	client.Register("schedule-test", scheduleTest)
	jobs := client.Command("cron-format", "schedule-test")
	jobs.EveryThreeHours()

	if !reflect.DeepEqual(jobs.Cron, []string{"0", "*/3", "*", "*", "*"}) {
		t.Fail()
	}
}

func TestEverySixHours(t *testing.T) {
	client.Register("schedule-test", scheduleTest)
	jobs := client.Command("cron-format", "schedule-test")
	jobs.EverySixHours()

	if !reflect.DeepEqual(jobs.Cron, []string{"0", "*/6", "*", "*", "*"}) {
		t.Fail()
	}
}

func TestEveryTwelveHours(t *testing.T) {
	client.Register("schedule-test", scheduleTest)
	jobs := client.Command("cron-format", "schedule-test")
	jobs.EveryTwelveHours()

	if !reflect.DeepEqual(jobs.Cron, []string{"0", "*/12", "*", "*", "*"}) {
		t.Fail()
	}
}

func TestDaily(t *testing.T) {
	client.Register("schedule-test", scheduleTest)
	jobs := client.Command("cron-format", "schedule-test")
	jobs.Daily()

	if !reflect.DeepEqual(jobs.Cron, []string{"0", "0", "*/1", "*", "*"}) {
		t.Fail()
	}
}

func TestDailyAt(t *testing.T) {
	client.Register("schedule-test", scheduleTest)
	jobs := client.Command("cron-format", "schedule-test")
	jobs.DailyAt("01:00")

	if !reflect.DeepEqual(jobs.Cron, []string{"0", "1", "*/1", "*", "*"}) {
		t.Fail()
	}
}

func TestWeekly(t *testing.T) {
	client.Register("schedule-test", scheduleTest)
	jobs := client.Command("cron-format", "schedule-test")
	jobs.Weekly()

	if !reflect.DeepEqual(jobs.Cron, []string{"0", "0", "*/7", "*", "*"}) {
		t.Fail()
	}
}

func TestWeeklyOn(t *testing.T) {
	client.Register("schedule-test", scheduleTest)
	jobs := client.Command("cron-format", "schedule-test")
	jobs.WeeklyOn("1:00")

	if !reflect.DeepEqual(jobs.Cron, []string{"0", "1", "*/7", "*", "*"}) {
		t.Fail()
	}
}

func TestMonthly(t *testing.T) {
	client.Register("schedule-test", scheduleTest)
	jobs := client.Command("cron-format", "schedule-test")
	jobs.Monthly()

	if !reflect.DeepEqual(jobs.Cron, []string{"0", "0", "1", "*/1", "*"}) {
		t.Fail()
	}
}

func TestMonthlyOn(t *testing.T) {
	client.Register("schedule-test", scheduleTest)
	jobs := client.Command("cron-format", "schedule-test")
	jobs.MonthlyOn("1:00")

	if !reflect.DeepEqual(jobs.Cron, []string{"0", "1", "1", "*/1", "*"}) {
		t.Fail()
	}
}

func TestQuarterly(t *testing.T) {
	client.Register("schedule-test", scheduleTest)
	jobs := client.Command("cron-format", "schedule-test")
	jobs.Quarterly()

	if !reflect.DeepEqual(jobs.Cron, []string{"0", "0", "1", "*/6", "*"}) {
		t.Fail()
	}
}

func TestYearly(t *testing.T) {
	client.Register("schedule-test", scheduleTest)
	jobs := client.Command("cron-format", "schedule-test")
	jobs.Yearly()

	if !reflect.DeepEqual(jobs.Cron, []string{"0", "0", "1", "*/12", "*"}) {
		t.Fail()
	}
}

func scheduleTest() error {
	return nil
}
