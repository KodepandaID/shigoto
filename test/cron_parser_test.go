package test

import (
	"testing"
	"time"

	cronparser "github.com/KodepandaID/shigoto/pkg/cron-parser"
)

func TestValidateMinute(t *testing.T) {
	str := []string{"*", "*/4", "0", "1", "59", "1-5"}
	parser := cronparser.New(&cronparser.Parser{
		Timezone: "Asia/Jakarta",
	})

	for _, row := range str {
		expr := []string{row, "*", "*", "*", "*"}
		if _, e := parser.Parse(expr); e != nil {
			t.Error(e)
			t.Fail()
		}
	}
}

func TestValidateFailMinute(t *testing.T) {
	str := []string{"60", "1,60", "0-60"}
	parser := cronparser.New(&cronparser.Parser{
		Timezone: "Asia/Jakarta",
	})

	var etotal int
	for _, row := range str {
		expr := []string{row, "*", "*", "*", "*"}
		if _, e := parser.Parse(expr); e != nil {
			etotal++
		}
	}

	if etotal < 3 {
		t.Error("This test should be fail")
		t.Fail()
	}
}

func TestValidateHour(t *testing.T) {
	str := []string{"*", "*/4", "1-23", "0"}
	parser := cronparser.New(&cronparser.Parser{
		Timezone: "Asia/Jakarta",
	})

	for _, row := range str {
		expr := []string{"*", row, "*", "*", "*"}
		if _, e := parser.Parse(expr); e != nil {
			t.Error(e)
			t.Fail()
		}
	}
}

func TestValidateFailHour(t *testing.T) {
	str := []string{"1/24", "24", "1,24"}
	parser := cronparser.New(&cronparser.Parser{
		Timezone: "Asia/Jakarta",
	})

	var etotal int
	for _, row := range str {
		expr := []string{"*", row, "*", "*", "*"}
		if _, e := parser.Parse(expr); e != nil {
			etotal++
		}
	}

	if etotal < 3 {
		t.Error("This test should be fail")
		t.Fail()
	}
}

func TestValidateDayMonth(t *testing.T) {
	str := []string{"*", "*/31", "1-31"}
	parser := cronparser.New(&cronparser.Parser{
		Timezone: "Asia/Jakarta",
	})

	for _, row := range str {
		expr := []string{"*", "*", row, "*", "*"}
		if _, e := parser.Parse(expr); e != nil {
			t.Error(e)
			t.Fail()
		}
	}
}

func TestValidateFailDayMonth(t *testing.T) {
	str := []string{"32", "0", "1-32"}
	parser := cronparser.New(&cronparser.Parser{
		Timezone: "Asia/Jakarta",
	})

	var etotal int
	for _, row := range str {
		expr := []string{"*", "*", row, "*", "*"}
		if _, e := parser.Parse(expr); e != nil {
			etotal++
		}
	}

	if etotal < 3 {
		t.Error("This test should be fail")
		t.Fail()
	}
}

func TestValidateMonth(t *testing.T) {
	str := []string{"*", "12", "1-12"}
	parser := cronparser.New(&cronparser.Parser{
		Timezone: "Asia/Jakarta",
	})

	for _, row := range str {
		expr := []string{"*", "*", "*", row, "*"}
		if _, e := parser.Parse(expr); e != nil {
			t.Error(e)
			t.Fail()
		}
	}
}

func TestValidateFailMonth(t *testing.T) {
	str := []string{"13", "0", "1-13"}
	parser := cronparser.New(&cronparser.Parser{
		Timezone: "Asia/Jakarta",
	})

	var etotal int
	for _, row := range str {
		expr := []string{"*", "*", "*", row, "*"}
		if _, e := parser.Parse(expr); e != nil {
			etotal++
		}
	}

	if etotal < 3 {
		t.Error("This test should be fail")
		t.Fail()
	}
}

func TestValidateWeekday(t *testing.T) {
	str := []string{"*", "0", "6"}
	parser := cronparser.New(&cronparser.Parser{
		Timezone: "Asia/Jakarta",
	})

	for _, row := range str {
		expr := []string{"*", "*", "*", "*", row}
		if _, e := parser.Parse(expr); e != nil {
			t.Error(e)
			t.Fail()
		}
	}
}

func TestValidateFailWeekday(t *testing.T) {
	str := []string{"7", "ABC", "1-7"}
	parser := cronparser.New(&cronparser.Parser{
		Timezone: "Asia/Jakarta",
	})

	var etotal int
	for _, row := range str {
		expr := []string{"*", "*", "*", "*", row}
		if _, e := parser.Parse(expr); e != nil {
			etotal++
		}
	}

	if etotal < 3 {
		t.Error("This test should be fail")
		t.Fail()
	}
}

func TestParseEveryMinute(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Date(2021, 5, 1, 1, 1, 0, 0, loc)

	valueTest := []string{
		"2021-05-01 01:02:00 +0700 WIB",
		"2021-05-01 01:03:00 +0700 WIB",
		"2021-05-01 01:04:00 +0700 WIB",
		"2021-05-01 01:05:00 +0700 WIB",
	}
	expr := []string{"*", "*", "*", "*", "*"}

	for _, val := range valueTest {
		parser := cronparser.New(&cronparser.Parser{
			Timezone: "Asia/Jakarta",
			SetTime:  now,
		})
		s, e := parser.Parse(expr)
		if e != nil {
			t.Error(e)
			t.Fail()
		}

		if s.Next.String() != val {
			t.Errorf("Value test not match: value is %s value should be %s", s.Next.String(), val)
			t.Fail()
		}
		now = s.Next
	}
}

func TestParseRangeMinute(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Date(2021, 5, 1, 1, 1, 0, 0, loc)

	valueTest := []string{
		"2021-05-01 01:02:00 +0700 WIB",
		"2021-05-01 02:01:00 +0700 WIB",
		"2021-05-01 02:02:00 +0700 WIB",
		"2021-05-01 03:01:00 +0700 WIB",
	}
	expr := []string{"1-2", "*", "*", "*", "*"}

	for _, val := range valueTest {
		parser := cronparser.New(&cronparser.Parser{
			Timezone: "Asia/Jakarta",
			SetTime:  now,
		})
		s, e := parser.Parse(expr)
		if e != nil {
			t.Error(e)
			t.Fail()
		}

		if s.Next.String() != val {
			t.Errorf("Value test not match: value is %s value should be %s", s.Next.String(), val)
			t.Fail()
		}
		now = s.Next
	}
}

func TestParseEveryStepMinute(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Date(2021, 5, 1, 1, 1, 0, 0, loc)

	valueTest := []string{
		"2021-05-01 01:15:00 +0700 WIB",
		"2021-05-01 01:30:00 +0700 WIB",
		"2021-05-01 01:45:00 +0700 WIB",
		"2021-05-01 02:00:00 +0700 WIB",
	}
	expr := []string{"*/15", "*", "*", "*", "*"}

	for _, val := range valueTest {
		parser := cronparser.New(&cronparser.Parser{
			Timezone: "Asia/Jakarta",
			SetTime:  now,
		})
		s, e := parser.Parse(expr)
		if e != nil {
			t.Error(e)
			t.Fail()
		}

		if s.Next.String() != val {
			t.Errorf("Value test not match: value is %s value should be %s", s.Next.String(), val)
			t.Fail()
		}
		now = s.Next
	}
}

func TestParseValueMinute(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Date(2021, 5, 1, 1, 1, 0, 0, loc)

	valueTest := []string{
		"2021-05-01 01:05:00 +0700 WIB",
		"2021-05-01 02:05:00 +0700 WIB",
		"2021-05-01 03:05:00 +0700 WIB",
		"2021-05-01 04:05:00 +0700 WIB",
	}
	expr := []string{"5", "*", "*", "*", "*"}

	for _, val := range valueTest {
		parser := cronparser.New(&cronparser.Parser{
			Timezone: "Asia/Jakarta",
			SetTime:  now,
		})
		s, e := parser.Parse(expr)
		if e != nil {
			t.Error(e)
			t.Fail()
		}

		if s.Next.String() != val {
			t.Errorf("Value test not match: value is %s value should be %s", s.Next.String(), val)
			t.Fail()
		}
		now = s.Next
	}
}

func TestParseEveryStepHourEveryMinute(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Date(2021, 5, 1, 1, 1, 0, 0, loc)

	valueTest := []string{
		"2021-05-01 01:02:00 +0700 WIB",
		"2021-05-01 01:03:00 +0700 WIB",
		"2021-05-01 01:04:00 +0700 WIB",
		"2021-05-01 01:05:00 +0700 WIB",
	}
	expr := []string{"*", "*/2", "*", "*", "*"}

	for _, val := range valueTest {
		parser := cronparser.New(&cronparser.Parser{
			Timezone: "Asia/Jakarta",
			SetTime:  now,
		})
		s, e := parser.Parse(expr)
		if e != nil {
			t.Error(e)
			t.Fail()
		}

		if s.Next.String() != val {
			t.Errorf("Value test not match: value is %s value should be %s", s.Next.String(), val)
			t.Fail()
		}
		now = s.Next
	}
}

func TestParseRangeHourEveryMinute(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Date(2021, 3, 17, 1, 57, 0, 0, loc)

	valueTest := []string{
		"2021-03-17 01:58:00 +0700 WIB",
		"2021-03-17 01:59:00 +0700 WIB",
		"2021-03-17 02:00:00 +0700 WIB",
		"2021-03-17 02:01:00 +0700 WIB",
		"2021-03-18 01:00:00 +0700 WIB",
	}
	expr := []string{"*", "1-2", "*", "*", "*"}

	for i, val := range valueTest {
		parser := cronparser.New(&cronparser.Parser{
			Timezone: "Asia/Jakarta",
			SetTime:  now,
		})
		s, e := parser.Parse(expr)
		if e != nil {
			t.Error(e)
			t.Fail()
		}

		if s.Next.String() != val {
			t.Errorf("Value test not match: value is %s value should be %s", s.Next.String(), val)
			t.Fail()
		}

		if i == 3 {
			now = time.Date(2021, 3, 17, 2, 59, 0, 0, loc)
		} else {
			now = s.Next
		}
	}
}

func TestParseEveryStepDayMonth(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	nowTimeCollection := []time.Time{
		time.Date(2021, 3, 2, 23, 58, 0, 0, loc),
		time.Date(2021, 3, 2, 23, 59, 0, 0, loc),
		time.Date(2021, 3, 4, 0, 0, 0, 0, loc),
		time.Date(2021, 3, 4, 23, 59, 0, 0, loc),
	}

	valueTest := []string{
		"2021-03-02 23:59:00 +0700 WIB",
		"2021-03-04 00:00:00 +0700 WIB",
		"2021-03-04 00:01:00 +0700 WIB",
		"2021-03-06 00:00:00 +0700 WIB",
	}

	expr := []string{"*", "*", "*/2", "*", "*"}

	for i, time := range nowTimeCollection {
		parser := cronparser.New(&cronparser.Parser{
			Timezone: "Asia/Jakarta",
			SetTime:  time,
		})
		s, e := parser.Parse(expr)
		if e != nil {
			t.Error(e)
			t.Fail()
		}

		if s.Next.String() != valueTest[i] {
			t.Errorf("Value test not match: value is %s value should be %s", s.Next.String(), valueTest[i])
			t.Fail()
		}
	}
}

func TestParseStepRangeDayMonth(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	nowTimeCollection := []time.Time{
		time.Date(2021, 3, 2, 23, 58, 0, 0, loc),
		time.Date(2021, 3, 2, 23, 59, 0, 0, loc),
		time.Date(2021, 4, 2, 0, 0, 0, 0, loc),
	}

	valueTest := []string{
		"2021-03-02 23:59:00 +0700 WIB",
		"2021-04-01 00:00:00 +0700 WIB",
		"2021-04-02 00:01:00 +0700 WIB",
	}

	expr := []string{"*", "*", "1-2", "*", "*"}

	for i, time := range nowTimeCollection {
		parser := cronparser.New(&cronparser.Parser{
			Timezone: "Asia/Jakarta",
			SetTime:  time,
		})
		s, e := parser.Parse(expr)
		if e != nil {
			t.Error(e)
			t.Fail()
		}

		if s.Next.String() != valueTest[i] {
			t.Errorf("Value test not match: value is %s value should be %s", s.Next.String(), valueTest[i])
			t.Fail()
		}
	}
}

func TestParseValueDayMonth(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	nowTimeCollection := []time.Time{
		time.Date(2021, 3, 2, 23, 58, 0, 0, loc),
		time.Date(2021, 3, 2, 23, 59, 0, 0, loc),
		time.Date(2021, 4, 2, 0, 0, 0, 0, loc),
	}

	valueTest := []string{
		"2021-03-02 23:59:00 +0700 WIB",
		"2021-04-02 00:00:00 +0700 WIB",
		"2021-04-02 00:01:00 +0700 WIB",
	}

	expr := []string{"*", "*", "2", "*", "*"}

	for i, time := range nowTimeCollection {
		parser := cronparser.New(&cronparser.Parser{
			Timezone: "Asia/Jakarta",
			SetTime:  time,
		})
		s, e := parser.Parse(expr)
		if e != nil {
			t.Error(e)
			t.Fail()
		}

		if s.Next.String() != valueTest[i] {
			t.Errorf("Value test not match: value is %s value should be %s", s.Next.String(), valueTest[i])
			t.Fail()
		}
	}
}

func TestParseEveryStepMonth(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	nowTimeCollection := []time.Time{
		time.Date(2021, 2, 1, 23, 59, 0, 0, loc),
		time.Date(2021, 2, 31, 23, 59, 0, 0, loc),
		time.Date(2021, 4, 1, 1, 0, 0, 0, loc),
	}

	valueTest := []string{
		"2021-02-02 00:00:00 +0700 WIB",
		"2021-04-01 00:00:00 +0700 WIB",
		"2021-04-01 01:01:00 +0700 WIB",
	}

	expr := []string{"*", "*", "*", "*/2", "*"}

	for i, time := range nowTimeCollection {
		parser := cronparser.New(&cronparser.Parser{
			Timezone: "Asia/Jakarta",
			SetTime:  time,
		})
		s, e := parser.Parse(expr)
		if e != nil {
			t.Error(e)
			t.Fail()
		}

		if s.Next.String() != valueTest[i] {
			t.Errorf("Value test not match: value is %s value should be %s", s.Next.String(), valueTest[i])
			t.Fail()
		}
	}
}

func TestParseValueMonth(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	nowTimeCollection := []time.Time{
		time.Date(2021, 3, 31, 23, 59, 0, 0, loc),
		time.Date(2021, 5, 31, 23, 58, 0, 0, loc),
		time.Date(2021, 5, 31, 23, 59, 0, 0, loc),
	}

	valueTest := []string{
		"2021-05-01 00:00:00 +0700 WIB",
		"2021-05-31 23:59:00 +0700 WIB",
		"2022-05-01 00:00:00 +0700 WIB",
	}

	expr := []string{"*", "*", "*", "5", "*"}

	for i, time := range nowTimeCollection {
		parser := cronparser.New(&cronparser.Parser{
			Timezone: "Asia/Jakarta",
			SetTime:  time,
		})
		s, e := parser.Parse(expr)
		if e != nil {
			t.Error(e)
			t.Fail()
		}

		if s.Next.String() != valueTest[i] {
			t.Errorf("Value test not match: value is %s value should be %s", s.Next.String(), valueTest[i])
			t.Fail()
		}
	}
}

func TestParseStepRangeMonth(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	nowTimeCollection := []time.Time{
		time.Date(2021, 1, 31, 23, 59, 0, 0, loc),
		time.Date(2021, 2, 1, 23, 58, 0, 0, loc),
		time.Date(2021, 2, 31, 23, 59, 0, 0, loc),
	}

	valueTest := []string{
		"2021-02-01 00:00:00 +0700 WIB",
		"2021-02-01 23:59:00 +0700 WIB",
		"2022-01-01 00:00:00 +0700 WIB",
	}

	expr := []string{"*", "*", "*", "1-2", "*"}

	for i, time := range nowTimeCollection {
		parser := cronparser.New(&cronparser.Parser{
			Timezone: "Asia/Jakarta",
			SetTime:  time,
		})
		s, e := parser.Parse(expr)
		if e != nil {
			t.Error(e)
			t.Fail()
		}

		if s.Next.String() != valueTest[i] {
			t.Errorf("Value test not match: value is %s value should be %s", s.Next.String(), valueTest[i])
			t.Fail()
		}
	}
}

func TestParseEveryStepWeek(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	nowTimeCollection := []time.Time{
		time.Date(2021, 3, 20, 5, 4, 0, 0, loc),
		time.Date(2021, 3, 21, 5, 4, 0, 0, loc),
	}

	valueTest := []string{
		"2021-03-21 05:04:00 +0700 WIB",
		"2021-03-23 05:04:00 +0700 WIB",
	}

	expr := []string{"4", "5", "*", "*", "*/2"}

	for i, time := range nowTimeCollection {
		parser := cronparser.New(&cronparser.Parser{
			Timezone: "Asia/Jakarta",
			SetTime:  time,
		})
		s, e := parser.Parse(expr)
		if e != nil {
			t.Error(e)
			t.Fail()
		}

		if s.Next.String() != valueTest[i] {
			t.Errorf("Value test not match: value is %s value should be %s", s.Next.String(), valueTest[i])
			t.Fail()
		}
	}
}

func TestParseStepRangeWeek(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	nowTimeCollection := []time.Time{
		time.Date(2021, 3, 22, 5, 4, 0, 0, loc),
		time.Date(2021, 3, 23, 5, 4, 0, 0, loc),
		time.Date(2021, 3, 30, 5, 4, 0, 0, loc),
	}

	valueTest := []string{
		"2021-03-23 05:04:00 +0700 WIB",
		"2021-03-29 05:04:00 +0700 WIB",
		"2021-04-05 05:04:00 +0700 WIB",
	}

	expr := []string{"4", "5", "*", "*", "1-2"}

	for i, time := range nowTimeCollection {
		parser := cronparser.New(&cronparser.Parser{
			Timezone: "Asia/Jakarta",
			SetTime:  time,
		})
		s, e := parser.Parse(expr)
		if e != nil {
			t.Error(e)
			t.Fail()
		}

		if s.Next.String() != valueTest[i] {
			t.Errorf("Value test not match: value is %s value should be %s", s.Next.String(), valueTest[i])
			t.Fail()
		}
	}
}

func TestParseValueWeek(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	nowTimeCollection := []time.Time{
		time.Date(2021, 3, 2, 23, 59, 0, 0, loc),
		time.Date(2021, 3, 2, 23, 59, 0, 0, loc),
		time.Date(2021, 3, 2, 23, 59, 0, 0, loc),
	}

	valueTest := []string{
		"2021-03-07 00:00:00 +0700 WIB",
		"2021-03-08 00:00:00 +0700 WIB",
		"2021-03-06 00:00:00 +0700 WIB",
	}

	exprs := [][]string{
		{"*", "*", "*", "*", "0"},
		{"*", "*", "*", "*", "1"},
		{"*", "*", "*", "*", "6"},
	}

	for i, expr := range exprs {
		parser := cronparser.New(&cronparser.Parser{
			Timezone: "Asia/Jakarta",
			SetTime:  nowTimeCollection[i],
		})
		s, e := parser.Parse(expr)
		if e != nil {
			t.Error(e)
			t.Fail()
		}

		if s.Next.String() != valueTest[i] {
			t.Errorf("Value test not match: value is %s value should be %s", s.Next.String(), valueTest[i])
			t.Fail()
		}
	}
}
