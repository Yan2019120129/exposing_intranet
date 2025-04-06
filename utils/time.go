package utils

import (
	"math/rand"
	"time"
)

// Time 随机
type Time struct {
	loc *time.Location
}

func NewTime() *Time {
	return &Time{}
}

// Intn 随机数
func (t *Time) Intn(min, max int) int {
	if min == max {
		return min
	}
	rand.NewSource(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func (t *Time) SetLocation(loc *time.Location) *Time {
	t.loc = loc
	return t
}

// RandTime 随机时间
func (t *Time) RandTime() time.Time {
	nowTime := time.Now()

	year := t.Intn(2000, nowTime.Year())

	month := time.Month(t.Intn(0, 12)) + 1

	nextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)

	lastDayOfMonth := nextMonth.AddDate(0, 0, -1)

	return t.randTime(year, month, t.Intn(0, lastDayOfMonth.Day()), t.Intn(0, 23), t.Intn(0, 60), t.Intn(0, 60), t.Intn(0, 999))
}

// RandToYearTime 随机当年时间
func (t *Time) RandToYearTime() time.Time {
	return t.randTime(0, time.Month(t.Intn(0, 12)), t.Intn(0, 23), t.Intn(0, 23), t.Intn(0, 60), t.Intn(0, 60), t.Intn(0, 999))
}

// RandToMonthTime 随机当月时间
func (t *Time) RandToMonthTime() time.Time {
	return t.randTime(0, 0, t.Intn(0, 12), t.Intn(0, 23), t.Intn(0, 60), t.Intn(0, 60), t.Intn(0, 999))
}

// RandToDayTime 随机当天时间
func (t *Time) RandToDayTime() time.Time {
	return t.randTime(0, 0, 0, t.Intn(0, 23), t.Intn(0, 60), t.Intn(0, 60), t.Intn(0, 999))
}

// RandToHourTime 随机当前小时的随机数
func (t *Time) RandToHourTime() time.Time {
	return t.randTime(0, 0, 0, 0, t.Intn(0, 60), t.Intn(0, 60), t.Intn(0, 999))
}

// RandToMinuteTime 随机当前分钟的随机数
func (t *Time) RandToMinuteTime() time.Time {
	return t.randTime(0, 0, 0, 0, 0, t.Intn(0, 60), t.Intn(0, 999))
}

// RandTime 随机时间
func (t *Time) randTime(year int, month time.Month, day, hour, minute, second, nanosecond int) time.Time {
	nowTime := time.Now()
	if year == 0 {
		year = nowTime.Year()
	}
	if month == 0 {
		month = nowTime.Month()
	}
	if day == 0 {
		day = nowTime.Day()
	}
	if hour == 0 {
		hour = nowTime.Hour()
	}
	if minute == 0 {
		minute = nowTime.Minute()
	}
	if second == 0 {
		second = nowTime.Second()
	}
	if nanosecond == 0 {
		nanosecond = nowTime.Nanosecond()
	}
	if t.loc == nil {
		t.loc = time.Local
	}
	return time.Date(year, month, day, hour, minute, second, nanosecond, t.loc)
}

// RandCustomizeTime 生成自定义随机时间
func (t *Time) RandCustomizeTime(year int, month time.Month, day, hour, minute, second, nanosecond int) time.Time {
	nowTime := time.Now()
	if year == 0 {
		year = t.Intn(2000, nowTime.Year())
	}
	if month == 0 {
		month = time.Month(t.Intn(1, 12))
	}
	if day == 0 {
		day = t.Intn(0, t.GetMonthDayNumb(year, int(month)))
	}
	if hour == 0 {
		hour = t.Intn(0, 24)
	}
	if minute == 0 {
		minute = t.Intn(0, 60)
	}
	if second == 0 {
		second = t.Intn(0, 60)
	}
	if nanosecond == 0 {
		nanosecond = t.Intn(0, 99999)
	}
	return t.randTime(year, month, day, hour, minute, second, nanosecond)
}

// GetMonthDayNumb 获取指定月份的天数
func (t *Time) GetMonthDayNumb(year, month int) int {
	nextMonth := time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, t.loc)
	lastDayOfMonth := nextMonth.AddDate(0, 0, -1)
	return lastDayOfMonth.Day()
}
