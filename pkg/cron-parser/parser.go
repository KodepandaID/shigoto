package cronparser

import (
	"math"
	"strconv"
	"strings"
	"time"
)

type Parser struct {
	loc         *time.Location
	Timezone    string
	SetTime     time.Time
	currentTime time.Time
}

// Schedule describes a job's duty cycle.
type Schedule struct {
	// Next returns the next activation time, later than the given time.
	// Next is invoked initially, and then each time the job is run.
	Next time.Time
}

type cronDirective struct {
	expr     string
	kind     int
	first    int
	last     int
	step     int
	stepSpan []int
}

// New to create a new instance
func New(p *Parser) Parser {
	return Parser{
		Timezone:    p.Timezone,
		currentTime: p.SetTime,
	}
}

// Parse to parsing cron format
func (p *Parser) Parse(expr []string) (Schedule, error) {
	var s Schedule
	if e := validate(expr); e != nil {
		return s, e
	}
	p.loc, _ = time.LoadLocation(p.Timezone)

	if p.Timezone != "" {
		if p.currentTime.IsZero() {
			p.currentTime = time.Now().Local().In(p.loc)
		}
	}

	d, e := p.exprParse(expr)
	if e != nil {
		return s, e
	}

	p.calculateMinute(d, p.loc)
	p.calculateHour(d, p.loc)
	p.calculateDay(d, p.loc)
	p.calculateMonth(d, p.loc)
	p.calculateWeek(d, p.loc)

	s.Next = p.currentTime

	return s, nil
}

func (p *Parser) calculateMinute(d []cronDirective, loc *time.Location) {
	m := p.currentTime.Minute()
	if d[0].kind == every && m <= d[0].last || d[0].kind != every && m <= d[0].last ||
		d[0].kind != one && d[1].kind != one && d[3].kind != one || d[3].kind == one && int(p.currentTime.Month()) == d[3].first {
		p.calculate(d[0], time.Minute, m, 60, 1)
	}
}

func (p *Parser) calculateHour(d []cronDirective, loc *time.Location) {
	h := p.currentTime.Hour()
	if d[1].kind == every && p.currentTime.Minute() > d[0].last || d[0].kind != every && p.currentTime.Minute() > d[0].last {
		p.currentTime = p.reset(p.currentTime, d[0].first, "minute", p.loc)
		p.calculate(d[1], time.Hour, h, 24, 1)
	} else if h > d[1].last {
		p.currentTime = p.reset(p.currentTime, d[1].first, "hour", p.loc)
		p.currentTime = p.currentTime.AddDate(0, 0, 1)
	}
}

func (p *Parser) calculateDay(d []cronDirective, loc *time.Location) {
	maxDay := daySpec(p.currentTime.Year())[p.currentTime.Month()-1]
	day := p.currentTime.Day()

	var nextstep int
	if len(d[2].stepSpan) > 0 {
		nextstep, _ = nextStep(p.currentTime.Day(), d[2].stepSpan)
	}

	if d[3].kind != one && d[2].kind == one && p.currentTime.Hour() == d[1].first && p.currentTime.Minute() == d[0].first ||
		d[3].kind != one && d[2].kind == everyStep && day < nextstep && p.currentTime.Hour() == d[1].first && p.currentTime.Minute() == d[0].first ||
		d[3].kind != one && d[2].kind == stepRange && day > d[2].last && p.currentTime.Hour() == d[1].first && p.currentTime.Minute() == d[0].first {
		p.currentTime = p.reset(p.currentTime, d[0].first, "minute", p.loc)
		p.currentTime = p.reset(p.currentTime, d[1].first, "hour", p.loc)
		p.calculate(d[2], time.Hour*24, day, maxDay, 1)
	}
}

func (p *Parser) calculateMonth(d []cronDirective, loc *time.Location) {
	m := int(p.currentTime.Month())

	var nextstep int
	var idx int
	if len(d[3].stepSpan) > 0 {
		nextstep, idx = nextStep(int(p.currentTime.Month()), d[3].stepSpan)
	}

	if d[3].kind == one && int(p.currentTime.Month()) != d[3].first && d[4].kind != one || d[3].kind == stepRange && int(p.currentTime.Month()) > d[3].last && d[4].kind != one {
		p.currentTime = p.reset(p.currentTime, d[0].first, "minute", p.loc)
		p.currentTime = p.reset(p.currentTime, d[1].first, "hour", p.loc)
		p.currentTime = p.reset(p.currentTime, d[2].first, "dayMonth", p.loc)

		if m >= d[3].first {
			p.currentTime = p.currentTime.AddDate(0, 12-m+d[3].first, 0)
		} else {
			p.currentTime = p.currentTime.AddDate(0, d[3].first-m, 0)
		}
	} else if d[3].kind == everyStep && int(p.currentTime.Month()) < nextstep && int(p.currentTime.Month()) > d[3].stepSpan[idx-1] &&
		p.currentTime.Minute() == d[0].first && p.currentTime.Hour() == d[1].first {
		p.currentTime = p.reset(p.currentTime, d[0].first, "minute", p.loc)
		p.currentTime = p.reset(p.currentTime, d[1].first, "hour", p.loc)
		p.currentTime = p.reset(p.currentTime, d[2].first, "dayMonth", p.loc)

		p.currentTime = p.currentTime.AddDate(0, nextstep-m, 0)
	}
}

func (p *Parser) calculateWeek(d []cronDirective, loc *time.Location) {
	if d[4].kind == one {
		nextweek := weekType[d[4].first]
		for p.currentTime.Weekday() != nextweek {
			p.currentTime = p.currentTime.AddDate(0, 0, 1)
		}
		if p.currentTime.Weekday() == nextweek {
			p.currentTime = p.reset(p.currentTime, d[0].first, "minute", p.loc)
			p.currentTime = p.reset(p.currentTime, d[1].first, "hour", p.loc)
		}
	} else if d[4].kind == stepRange && p.currentTime.Hour() == d[1].first && p.currentTime.Minute() == d[0].first {
		for p.currentTime.Weekday() != weekType[d[4].first] {
			if p.currentTime.Weekday() == weekType[d[4].last] {
				break
			}
			p.currentTime = p.currentTime.AddDate(0, 0, 1)
		}
		if p.currentTime.Weekday() == weekType[d[4].first] || p.currentTime.Weekday() == weekType[d[4].last] {
			p.currentTime = p.reset(p.currentTime, d[0].first, "minute", p.loc)
			p.currentTime = p.reset(p.currentTime, d[1].first, "hour", p.loc)
		}
	} else if d[4].kind == everyStep && p.currentTime.Minute() == d[0].first && p.currentTime.Hour() == d[1].first {
		p.currentTime = p.reset(p.currentTime, d[0].first, "minute", p.loc)
		p.currentTime = p.reset(p.currentTime, d[1].first, "hour", p.loc)

		var nextstep int
		var idx int
		if len(d[4].stepSpan) > 0 {
			nextstep, idx = nextStep(p.currentTime.Day(), d[4].stepSpan)
		}

		if p.currentTime.Day() != d[4].stepSpan[idx-1] {
			p.currentTime = p.currentTime.AddDate(0, 0, nextstep-p.currentTime.Day())
		}
	}

}

func (p *Parser) exprParse(expr []string) (directive []cronDirective, e error) {
	lastTimeCollection := []int{59, 23, daySpec(p.currentTime.Year())[11], 12, 6}
	for i, val := range expr {
		sval := strings.ToLower(val)
		dtmp := cronDirective{
			expr: sval,
			kind: 0,
		}
		pattern := regexTimeCollection[i]
		lastTime := lastTimeCollection[i]

		// `*`
		if makeLayoutRegexp(layoutWildcard, pattern).MatchString(sval) {
			dtmp.kind = every
			dtmp.step = 1
			if i < 2 {
				dtmp.first = 0
			} else {
				dtmp.first = 1
			}
			dtmp.last = lastTime

			directive = append(directive, dtmp)
			continue
		}
		// `1`
		if makeLayoutRegexp(layoutValue, pattern).MatchString(sval) {
			dtmp.kind = one

			d, e := strconv.Atoi(sval)
			if e != nil {
				dtmp.step = 0
			}
			dtmp.step = d
			dtmp.first = d
			dtmp.last = d

			directive = append(directive, dtmp)
			continue
		}
		// `1-2`
		pairs := makeLayoutRegexp(layoutRange, pattern).FindStringSubmatchIndex(sval)
		if len(pairs) > 0 {
			dtmp.kind = stepRange
			first, _ := strconv.Atoi(sval[pairs[2]:pairs[3]])
			dtmp.first = first

			second, _ := strconv.Atoi(sval[pairs[4]:pairs[5]])
			dtmp.last = second
			dtmp.step = 1

			directive = append(directive, dtmp)
			continue
		}

		// `*/2`
		pairs = makeLayoutRegexp(layoutWildcardAndInterval, pattern).FindStringSubmatchIndex(sval)
		if len(pairs) > 0 {
			var lspan int
			j := 0
			step := 0

			dtmp.kind = everyStep
			d, _ := strconv.Atoi(sval[pairs[2]:pairs[3]])

			switch exprType[i] {
			case "minute":
				lspan = int(math.Round(float64(60 / d)))
			case "hour":
				lspan = int(math.Round(float64(24 / d)))
			case "dayMonth":
				maxDay := daySpec(p.currentTime.Year())[p.currentTime.Month()-1]
				lspan = int(math.Round(float64(maxDay / d)))
			case "month":
				lspan = int(math.Round(float64(12 / d)))
			case "weekday":
				maxDay := daySpec(p.currentTime.Year())[p.currentTime.Month()-1]
				lspan = int(math.Round(float64(maxDay / d)))
				step = 1
			}

			for j <= lspan {
				dtmp.stepSpan = append(dtmp.stepSpan, step)
				j++
				step += d
			}

			dtmp.last = dtmp.stepSpan[len(dtmp.stepSpan)-1]
			directive = append(directive, dtmp)
			continue
		}
	}

	return directive, e
}

func (p *Parser) reset(t time.Time, first int, kind string, loc *time.Location) time.Time {
	switch kind {
	case "minute":
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), first, 0, 0, loc)
	case "hour":
		return time.Date(t.Year(), t.Month(), t.Day(), first, t.Minute(), 0, 0, loc)
	case "dayMonth":
		return time.Date(t.Year(), t.Month(), first, t.Hour(), t.Minute(), 0, 0, loc)
	case "month":
		return time.Date(t.Year(), time.Month(first), t.Day(), t.Hour(), t.Minute(), 0, 0, loc)
	}

	return t
}

func (p *Parser) calculate(d cronDirective, td time.Duration, tnow, max, seconds int) bool {
	var next bool
	var diff int
	switch d.kind {
	case one:
		if tnow < d.step {
			diff = d.step - tnow
		} else {
			diff = max - tnow + d.step
		}
	case stepRange:
		if tnow < d.first {
			diff = d.first - tnow
		} else if tnow < d.last {
			diff = 1
		} else if tnow >= d.last {
			diff = (max - tnow) + d.first
		}
	case everyStep:
		step, _ := nextStep(tnow, d.stepSpan)
		diff = step - tnow
	case every:
		diff = d.step * seconds
	}

	p.currentTime = p.currentTime.Add(td * time.Duration(diff))

	return next
}
