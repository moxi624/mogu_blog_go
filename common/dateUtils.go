package common

import (
	"strings"
	"time"
)

/**
 *
 * @author  镜湖老杨
 * @date  2020/12/23 3:46 下午
 * @version 1.0
 */

type dateUtils struct{}

func (dateUtils) GetNowTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func (dateUtils) GetDateStr(beforeDate time.Time, timeSecond int64) string {
	time := beforeDate.Add(time.Duration(timeSecond) * time.Second)
	return time.Format("2006-01-02 15:04:05")
}

func (dateUtils) GetToDayEndTime() string {
	t := " 23:59:59"
	return time.Now().Format("2006-01-02") + t
}

func (dateUtils) GetDate(date string, day time.Duration) string {
	t, _ := time.Parse("2006-01-02 15:04:05", date)
	t1 := t.Add(time.Hour * 24 * day)
	return t1.Format("2006-01-02 15:04:05")
}

func (dateUtils) GetDaysByN(intervals time.Duration, formatStr string) []string {
	var pastDaysList []string
	t, _ := time.Parse(formatStr, time.Now().Format(formatStr))
	for i := intervals - 1; i >= 0; i-- {
		t1 := t.Add(-time.Hour * 24 * i)
		pastDaysList = append(pastDaysList, t1.Format(formatStr))
	}
	return pastDaysList
}

func (dateUtils) GetToDayStartTime() string {
	t := " 00:00:00"
	return time.Now().Format("2006-01-02") + t
}

func (dateUtils) GetDayBetweenDates(begin, end string) []string {
	var lDate []string
	s := strings.Split(begin, " ")
	lDate = append(lDate, s[0])
	calBegin, _ := time.Parse("2006-01-02 15:04:05", begin)
	calEnd, _ := time.Parse("2006-01-02 15:04:05", end)
	for {
		if !calEnd.After(calBegin) {
			break
		}
		calBegin = calBegin.AddDate(0, 0, 1)
		lDate = append(lDate, calBegin.Format("2006-01-02"))
	}
	return lDate
}

func (dateUtils) Str2Date(dataString string) string {
	dataString = strings.Split(dataString, "(中国标准时间)")[0]
	dataString = strings.SplitN(dataString, " ", 2)[1]
	dataString = strings.Split(dataString, " GMT")[0]
	date, err := time.Parse("Jan 02 2006 15:04:05", dataString)
	if err != nil {
		panic(err)
	}
	return date.Format("2006-01-02 15:04:05")
}

var DateUtils = &dateUtils{}
