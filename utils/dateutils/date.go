package dateutils

import "time"

const (
	dateTimeFormat = "2006-01-02 15:04:05"
)

type DateTime time.Time

func GetDateTimeNow() time.Time {
	return time.Now().UTC()
}

func GetDateTimeNowFormat() string {
	return GetDateTimeNow().Format(dateTimeFormat)
}
