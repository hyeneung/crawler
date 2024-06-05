package utils

import (
	"errors"
	"time"
)

func UnixTime2Time(unixTime int64) time.Time {
	return time.Unix(unixTime, 0).UTC()
}

func Str2time(strTime string) time.Time {
	timeFormats := [...]string{
		time.RFC1123,
		time.RFC1123Z,
		time.RFC822Z,
		time.RFC850,
		time.RFC822,
		time.RFC3339,
	}
	// RFC1123 : Thu, 02 May 2024 08:00:00 GMT
	// https://go.dev/src/time/format.go
	for _, format := range timeFormats {
		t, err := time.Parse(format, strTime)
		if err == nil {
			return t.UTC()
		}
	}

	CheckErr(errors.New("time format error"), SlogLogger)
	return time.Time{}
}

func Str2UnixTime(strTime string) int64 {
	t := Str2time(strTime)
	return t.Unix()
}
