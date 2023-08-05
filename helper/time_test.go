package helper

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	parseTime, err := ParseTime("")
	assert.NotNil(t, err)
	assert.Nil(t, parseTime)
}

func TestGetIntervalTime(t *testing.T) {
	req := &GetIntervalTimeRequest{
		Type: Day,
		Num:  0,
	}
	req.Time, _ = time.ParseInLocation(time.DateTime, "2023-08-05 00:00:00", time.Local)
	// test day
	startAt, _ := time.ParseInLocation(time.DateTime, "2023-08-05 00:00:00", time.Local)
	endAt, _ := time.ParseInLocation(time.DateTime, "2023-08-05 23:59:59", time.Local)
	resp := GetIntervalTime(req)
	assert.Equal(t, startAt, resp.StartAt)
	assert.Equal(t, endAt, resp.EndAt)
	// test week
	req.Type = Week
	startAt, _ = time.ParseInLocation(time.DateTime, "2023-07-31 00:00:00", time.Local)
	endAt, _ = time.ParseInLocation(time.DateTime, "2023-08-06 23:59:59", time.Local)
	resp = GetIntervalTime(req)
	assert.Equal(t, startAt, resp.StartAt)
	assert.Equal(t, endAt, resp.EndAt)
}
