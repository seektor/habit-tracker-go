package habit

import (
	"fmt"
)

type TotalTime struct {
	Days    int16
	Hours   int8
	Minutes int8
}

func (t TotalTime) Stringify() string {
	text := ""

	if t.Days > 0 {
		if t.Days == 1 {
			text += "1 Day"
		} else {
			text += fmt.Sprintf("%d Days", t.Days)
		}
	}

	if t.Hours > 0 {
		if len(text) > 0 {
			text += " "
		}

		if t.Hours == 1 {
			text += "1 Hour"
		} else {
			text += fmt.Sprintf("%d Hours", t.Hours)
		}
	}

	if t.Minutes > 0 {
		if len(text) > 0 {
			text += " "
		}

		if t.Minutes == 1 {
			text += "1 Minute"
		} else {
			text += fmt.Sprintf("%d Minutes", t.Minutes)
		}
	}

	if text == "" {
		return "-"
	}

	return text
}

func (t *TotalTime) Add(minutes int16) {
	totalMinutes := int16(t.Minutes) + minutes
	newMinutes := int8(totalMinutes % 60)
	extraHours := int8(totalMinutes / 60)

	totalHours := t.Hours + extraHours
	newHours := totalHours % 24
	extraDays := totalHours / 24

	newDays := t.Days + int16(extraDays)

	t.Days = newDays
	t.Hours = newHours
	t.Minutes = newMinutes
}

func (t *TotalTime) subtractDays(days int16) {
	newDays := t.Days - days

	if newDays < 0 {
		t.Days = 0
		t.Hours = 0
		t.Minutes = 0
	} else {
		t.Days = newDays
	}

}

func (t *TotalTime) subtractHours(hours int8) {
	totalHours := t.Hours - hours

	if totalHours >= 0 {
		t.Hours = totalHours
	} else {
		clampedHours := -totalHours % 24
		t.Hours = (24 - clampedHours) % 24

		daysToSubtract := int16(-totalHours/25) + 1
		t.subtractDays(daysToSubtract)
	}
}

func (t *TotalTime) subtractMinutes(minutes int16) {
	totalMinutes := int16(t.Minutes) - minutes

	if totalMinutes >= 0 {
		t.Minutes = int8(totalMinutes)
	} else {
		clampedMinutes := -totalMinutes % 60
		t.Minutes = int8((60 - clampedMinutes) % 60)

		hoursToSubtract := int8(-totalMinutes/61) + 1
		t.subtractHours(hoursToSubtract)
	}
}

func (t *TotalTime) Subtract(minutes int16) {
	t.subtractMinutes(minutes)
}
