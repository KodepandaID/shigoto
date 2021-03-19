package cronparser

import (
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	one                       = 1
	everyStep                 = 2
	stepRange                 = 3
	every                     = 4
	exprType                  = []string{"minute", "hour", "dayMonth", "month", "weekday"}
	regexTimeCollection       = []string{minuteValue, hourValue, dayMonthValue, monthValue, weekdayValue}
	weekType                  = []time.Weekday{time.Sunday, time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday}
	layoutWildcard            = `^[*]$`
	layoutValue               = `^(%value%)$`
	layoutRange               = `^(%value%)-(%value%)$`
	layoutWildcardAndInterval = `^\*/(\d+)$`
	layoutValueAndInterval    = `^(%value%)/(\d+)$`
	layoutRangeAndInterval    = `^(%value%)-(%value%)/(\d+)$`
	layoutDowOfSpecificWeek   = `^(%value%)#([1-5])$`
	fieldFinder               = regexp.MustCompile(`\S+`)
	entryFinder               = regexp.MustCompile(`[^,]+`)
	layoutRegexp              = make(map[string]*regexp.Regexp)
	layoutRegexpLock          sync.Mutex
)

var normalDaySpec = []int{
	31, // January
	28, // February
	31, // March
	30, // April
	31, // May
	30, // June
	31, // July
	31, // August
	30, // September
	31, // October
	30, // November
	31, // December
}

var leapDaySpec = []int{
	31, // January
	29, // February
	31, // March
	30, // April
	31, // May
	30, // June
	31, // July
	31, // August
	30, // September
	31, // October
	30, // November
	31, // December
}

func nextStep(val int, nums []int) (next int, index int) {
	for i := 0; i < len(nums); i++ {
		if i != len(nums)-1 {
			if val > nums[i] && val < nums[i+1] || val == nums[i] {
				return nums[i+1], i + 1
			}
		} else {
			return nums[len(nums)-1], len(nums) - 1
		}
	}

	return nums[1], 1
}

func makeLayoutRegexp(layout, value string) *regexp.Regexp {
	layoutRegexpLock.Lock()
	defer layoutRegexpLock.Unlock()

	layout = strings.Replace(layout, `%value%`, value, -1)
	re := layoutRegexp[layout]
	if re == nil {
		re = regexp.MustCompile(layout)
		layoutRegexp[layout] = re
	}

	return re
}

func daySpec(year int) []int {
	if year%400 == 0 || (year%4 == 0 && year%100 != 0) {
		return leapDaySpec
	}

	return normalDaySpec
}
