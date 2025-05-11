package utils

import "time"

func getBeginningOfDayDate(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func GetDaysDiff(from time.Time, to time.Time) int32 {
	fromBeginning := getBeginningOfDayDate(from)
	toBeginning := getBeginningOfDayDate(to)

	return int32(toBeginning.Sub(fromBeginning).Hours() / 24)
}
