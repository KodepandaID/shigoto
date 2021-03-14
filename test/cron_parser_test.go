package test

import (
	"fmt"
	"testing"

	cronparser "github.com/KodepandaID/shigoto/pkg/cron-parser"
)

func TestValidateMinute(t *testing.T) {
	str := []string{"*", "1/4", "0", "1", "59", "1,5", "1-5"}
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
	str := []string{"*", "1/4", "1-23", "1,2", "0"}
	parser := cronparser.New(&cronparser.Parser{
		Timezone: "Asia/Jakarta",
	})

	for _, row := range str {
		expr := []string{"*", row, "*", "*", "*"}
		if _, e := parser.Parse(expr); e != nil {
			fmt.Println(row)
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
	str := []string{"*", "1/31", "1-31", "1,2"}
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
	str := []string{"*", "0", "6", "SUN"}
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
