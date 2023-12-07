package helper

import "time"

// ParseTime is a function that parses a string into a time.Time object.
// It takes a string representing a time in the format "2006-01-02 15:04:05" as a parameter,
// and returns a pointer to a time.Time object and an error.
// If the string cannot be parsed into a time.Time object, the function returns nil and the error.
func ParseTime(t string) (*time.Time, error) {
	res, err := time.Parse("2006-01-02 15:04:05", t)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// IntervalTimeType is a type that represents the type of a time interval.
type IntervalTimeType int

// Constants for the different types of time intervals.
const (
	Day IntervalTimeType = iota
	Week
	Month
	Year
)

// GetIntervalTimeRequest is a struct that represents a request to get the start and end times of a time interval.
// It contains the time that specifies the date, the type of the time interval, and the number that specifies the offset from the specified date.
type GetIntervalTimeRequest struct {
	// Time is the time that specifies the date.
	Time time.Time
	// Type is the type of the time interval.
	Type IntervalTimeType
	// Num is the number that specifies the offset from the specified date.
	Num int
}

// GetIntervalTimeResponse is a struct that represents the response to a GetIntervalTimeRequest.
// It contains the start and end times of the time interval.
type GetIntervalTimeResponse struct {
	// StartAt is the start time of the time interval.
	StartAt time.Time
	// EndAt is the end time of the time interval.
	EndAt time.Time
}

// GetIntervalTime is a function that gets the start and end times of a time interval.
// It takes a pointer to a GetIntervalTimeRequest struct as a parameter,
// and returns a pointer to a GetIntervalTimeResponse struct.
// The function calculates the start and end times based on the type and number in the GetIntervalTimeRequest,
// and sets them in the GetIntervalTimeResponse.
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
