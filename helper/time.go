package helper

import "time"

func ParseTime(t string) (*time.Time, error) {
	res, err := time.Parse("2006-01-02 15:04:05", t)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// WeekIntervalTime 获取某周的开始和结束时间,week为0本周,-1上周，1下周以此类推
func WeekIntervalTime(week int) (startTime, endTime time.Time) {
	now := time.Now()
	offset := int(time.Monday - now.Weekday())
	//周日做特殊判断 因为time.Monday = 0
	if offset > 0 {
		offset = -6
	}

	year, month, day := now.Date()
	thisWeek := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	startTime, _ = time.ParseInLocation("2006-01-02 15:04:05", thisWeek.AddDate(0, 0, offset+7*week).Format("2006-01-02")+" 00:00:00", time.Local)
	endTime, _ = time.ParseInLocation("2006-01-02 15:04:05", thisWeek.AddDate(0, 0, offset+6+7*week).Format("2006-01-02")+" 23:59:59", time.Local)

	return
}

// MonthIntervalTime 获取某月的开始和结束时间,m为0本月,-1上月，1下月以此类推
func MonthIntervalTime(m int) (startTime, endTime time.Time) {
	year, month, _ := time.Now().AddDate(0, m, 0).Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	startTime, _ = time.ParseInLocation("2006-01-02 15:04:05", thisMonth.Format("2006-01-02")+" 00:00:00", time.Local)
	endTime, _ = time.ParseInLocation("2006-01-02 15:04:05", thisMonth.AddDate(0, 1, -1).Format("2006-01-02")+" 23:59:59", time.Local)
	return
}

// YearIntervalTime 获取某年的开始和结束时间,y为0今年,-1去年，1明年以此类推
func YearIntervalTime(y int) (startTime, endTime time.Time) {
	year, _, _ := time.Now().AddDate(y, 0, 0).Date()
	thisYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
	startTime, _ = time.ParseInLocation("2006-01-02 15:04:05", thisYear.Format("2006-01-02")+" 00:00:00", time.Local)
	endTime, _ = time.ParseInLocation("2006-01-02 15:04:05", thisYear.AddDate(1, 0, -1).Format("2006-01-02")+" 23:59:59", time.Local)
	return
}
