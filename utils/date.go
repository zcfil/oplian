package utils

import "time"

const (
	YMDHMS          = "20060102150405"
	YearMonthDayHMS = "2006-01-02 15:04:05"
	YearMonthDayHM  = "2006-01-02 15:04"
	YearMonthDay    = "2006-01-02"
	Year            = "2006"
	YearMonth       = "200601"
	Month           = "01"
	Day             = "02"
)

func GetNowStr() string {
	return time.Now().Format(YearMonthDayHMS)
}

func DateToTimeStamp(date string, timeLayout string) (int64, error) {
	if timeLayout == "" {
		timeLayout = "2006-01-02 15:04:05"
	}
	loc, err := time.LoadLocation("Local") //获取时区
	tmp, err := time.ParseInLocation(timeLayout, date, loc)
	return tmp.Unix(), err //转化为时间戳 类型是int64
}

func StrToTime(data string) time.Time {

	loc, _ := time.LoadLocation("Local")
	t, _ := time.ParseInLocation(YearMonthDayHMS, data, loc)
	return t
}

func TimeToFormat(t time.Time, timeType string) string {

	timeFormat := ""
	switch timeType {
	case Year:
		timeFormat = Year
		break
	case Month:
		timeFormat = Month
		break
	case Day:
		timeFormat = Day
		break
	case YearMonth:
		timeFormat = YearMonth
		break
	case YearMonthDayHM:
		timeFormat = YearMonthDayHM
		break
	case YearMonthDay:
		timeFormat = YearMonthDay
		break
	case YMDHMS:
		timeFormat = YMDHMS
		break
	default:
		timeFormat = YearMonthDayHMS
	}
	return t.Format(timeFormat)

}

func TimeAddDay(day int) string {
	return TimeToFormat(time.Now().AddDate(0, 0, day), YearMonthDay)
}
