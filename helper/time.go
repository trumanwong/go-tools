package helper

import "time"

func ParseTime(t string) (*time.Time, error) {
	res, err := time.Parse("2006-01-02 15:04:05", t)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type IntervalTimeType int

const (
	Day IntervalTimeType = iota
	Week
	Month
	Year
)

type GetIntervalTimeRequest struct {
	// 指定日期
	Time time.Time
	Type IntervalTimeType
	// 如果类型为Day， 则Num为0指定日期的当天，-1指定日期的前一天，1指定日期的后一天以此类推
	// 如果类型为Week， 则Num为0指定日期的当周,-1指定日期的上周，1指定日期的下周以此类推
	// 如果类型为Month， 则Num为0指定日期的当月,-1指定日期的上月，1指定日期的下月以此类推
	// 如果类型为Year， 则Num为0指定日期的当年,-1指定日期的去年，1指定日期的明年以此类推
	Num int
}

type GetIntervalTimeResponse struct {
	StartAt time.Time
	EndAt   time.Time
}

// GetIntervalTime 获取某天/周/月/年的开始和结束时间,d为0今天,-1昨天，1明天以此类推
func GetIntervalTime(req *GetIntervalTimeRequest) (resp *GetIntervalTimeResponse) {
	resp = new(GetIntervalTimeResponse)
	var offset int
	switch req.Type {
	case Day:
		t := req.Time.AddDate(0, 0, req.Num)
		resp.StartAt, _ = time.ParseInLocation(time.DateTime, t.Format(time.DateOnly)+" 00:00:00", time.Local)
		resp.EndAt, _ = time.ParseInLocation(time.DateTime, t.Format(time.DateOnly)+" 23:59:59", time.Local)
	case Week:
		year, month, day := req.Time.Date()
		offset = int(time.Monday - req.Time.Weekday())
		//周日做特殊判断 因为time.Monday = 0
		if offset > 0 {
			offset = -6
		}
		thisWeek := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
		resp.StartAt, _ = time.ParseInLocation(time.DateTime, thisWeek.AddDate(0, 0, offset+7*req.Num).Format(time.DateOnly)+" 00:00:00", time.Local)
		resp.EndAt, _ = time.ParseInLocation(time.DateTime, thisWeek.AddDate(0, 0, offset+6+7*req.Num).Format(time.DateOnly)+" 23:59:59", time.Local)
	case Month:
		year, month, _ := req.Time.AddDate(0, req.Num, 0).Date()
		thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
		resp.StartAt, _ = time.ParseInLocation(time.DateTime, thisMonth.Format(time.DateOnly)+" 00:00:00", time.Local)
		resp.EndAt, _ = time.ParseInLocation(time.DateTime, thisMonth.AddDate(0, 1, -1).Format(time.DateOnly)+" 23:59:59", time.Local)
	case Year:
		year, _, _ := req.Time.AddDate(req.Num, 0, 0).Date()
		thisYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
		resp.StartAt, _ = time.ParseInLocation(time.DateTime, thisYear.Format(time.DateOnly)+" 00:00:00", time.Local)
		resp.EndAt, _ = time.ParseInLocation(time.DateTime, thisYear.AddDate(1, 0, -1).Format(time.DateOnly)+" 23:59:59", time.Local)
	}
	return
}
