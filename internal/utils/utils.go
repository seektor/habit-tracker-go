package utils

import (
	"fmt"
	"time"
)

const FileName = "habits_tracker.json"

var FgColors = struct {
	Reset  string
	Yellow string
	Green  string
	Red    string
	Blue   string
	Bold   string
}{
	Reset:  "\033[0m",
	Yellow: "\033[33m",
	Green:  "\033[32m",
	Red:    "\033[31m",
	Blue:   "\033[34m",
	Bold:   "\033[1m",
}

func getBeginningOfDayDate(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func GetDaysDiff(from time.Time, to time.Time) int32 {
	fromBeginning := getBeginningOfDayDate(from)
	toBeginning := getBeginningOfDayDate(to)

	return int32(toBeginning.Sub(fromBeginning).Hours() / 24)
}

func ColorString(color string, msg string) string {
	return fmt.Sprint(color + msg + FgColors.Reset)
}

func PrintlnError(msg string) {
	fmt.Println(FgColors.Red + "=== " + msg + " ===" + FgColors.Reset)
}

func PrintlnSuccess(msg string) {
	fmt.Println(FgColors.Green + "=== " + msg + " ===" + FgColors.Reset)
}

func PrintlnInfo(msg string) {
	fmt.Println("=== " + msg + " ===")
}
