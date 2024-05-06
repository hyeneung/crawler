package utils

import "time"

func UnixTime2Time(unixTime int64) time.Time {
	return time.Unix(unixTime, 0).UTC()
}

func Str2time(strTime string) time.Time {
	// RFC1123 : Thu, 02 May 2024 08:00:00 GMT
	t, err := time.Parse(time.RFC1123, strTime)
	CheckErr(err)
	return t.UTC()
}
