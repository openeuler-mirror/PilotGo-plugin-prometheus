package utils

import (
	"strconv"
	"time"
)

func UnixTimeToShanghai(unix string) (time.Time, error) {
	unixInt, err := strconv.ParseInt(unix, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	unixSecs := unixInt / 1000
	unixNsecs := (unixInt % 1000) * int64(time.Millisecond)

	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return time.Time{}, err
	}

	utcTime := time.Unix(unixSecs, unixNsecs)
	cstTime := utcTime.In(loc)
	return cstTime, nil
}

func ShanghaiTimeToUnixMillis(dateStr string) int64 {
	t, _ := time.Parse("2006-01-02", dateStr)

	utcTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).UTC()

	unixSecs := utcTime.Unix()
	unixNsecs := int64(utcTime.Nanosecond())

	unixMillis := unixSecs*1000 + unixNsecs/int64(time.Millisecond)

	return unixMillis
}
