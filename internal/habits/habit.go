package habits

import (
	"fmt"
	"time"
)

const MaxHabitNameLength int8 = 16
const MaxHabitTotalTime int16 = 16 * 60 // minutes
const HistoryLen int8 = 6

type Entry struct {
	CheckedSteps int8
	StepsCount   int8
	IsFrozen     bool
}
type Summary struct {
	TotalTime     TotalTime
	LongestStreak int16
	CurrentStreak int16
	History       [HistoryLen]Entry // History of last 6 days
}

type Habit struct {
	Name         string
	CreatedAt    time.Time
	StepsCount   int8
	StepMinutes  int16 // minutes
	CheckedSteps int8
	IsFrozen     bool
	Summary      Summary
}

func newHabit(name string, stepsCount int8, stepTime int16) Habit {
	habit := Habit{
		Name:         name,
		CreatedAt:    time.Now(),
		StepsCount:   stepsCount,
		StepMinutes:  stepTime,
		CheckedSteps: 0,
		IsFrozen:     false,
		Summary: Summary{
			TotalTime:     TotalTime{},
			LongestStreak: 0,
			CurrentStreak: 0,
			History:       [HistoryLen]Entry{},
		},
	}

	return habit
}

func (h *Habit) CheckStep() {
	if h.IsFrozen {
		return
	}

	h.CheckedSteps += 1
}

func (h *Habit) UncheckStep() {
	if h.IsFrozen {
		return
	}

	if h.CheckedSteps > 0 {
		h.CheckedSteps -= 1
	}
}

func validateStepData(stepsCount int8, stepTime int16) error {
	if stepsCount < 1 || stepTime < 1 {
		return fmt.Errorf("step count and Step time has to be a positive value")
	}

	if int16(stepsCount)*stepTime > MaxHabitTotalTime {
		return fmt.Errorf("max habit total time cannot exceed %d", MaxHabitTotalTime)
	}

	return nil
}

func (h *Habit) SetStepsCount(stepsCount int8) error {
	err := validateStepData(stepsCount, h.StepMinutes)

	if err != nil {
		return err
	}

	h.StepsCount = stepsCount

	return nil
}

func (h *Habit) SetStepMinutes(stepMinutes int16) error {
	err := validateStepData(h.StepsCount, stepMinutes)

	if err != nil {
		return err
	}

	h.StepMinutes = stepMinutes

	return nil
}

func (h *Habit) Freeze() {
	h.IsFrozen = true
	h.CheckedSteps = 0
}

func (h *Habit) Unfreeze() {
	h.IsFrozen = false
}

func (h *Habit) getCurrentEntry() Entry {
	if h.IsFrozen {
		return Entry{IsFrozen: true}
	} else {
		return Entry{
			CheckedSteps: h.CheckedSteps,
			StepsCount:   h.StepsCount,
		}
	}
}

func (h *Habit) updateStatistics() {
	if h.IsFrozen {
		return
	}

	h.Summary.TotalTime.Add(h.StepMinutes * int16(h.CheckedSteps))

	if h.CheckedSteps >= h.StepsCount {
		h.Summary.CurrentStreak += 1

		if h.Summary.CurrentStreak > h.Summary.LongestStreak {
			h.Summary.LongestStreak = h.Summary.CurrentStreak
		}
	} else {
		h.Summary.CurrentStreak = 0
	}
}

func (h *Habit) UpdateToPresent(daysDiff int32) {
	if daysDiff <= 0 {
		return
	}

	firstDayEntryIdx := int32(HistoryLen) - daysDiff

	if firstDayEntryIdx >= 0 {
		// Shift history and save the current entry when it is within the history range
		for i := range int(firstDayEntryIdx) {
			h.Summary.History[i] = h.Summary.History[i+int(daysDiff)]
		}

		h.Summary.History[firstDayEntryIdx] = h.getCurrentEntry()
	}

	// Update based on the values from the day of the last user activity
	h.updateStatistics()
	h.CheckedSteps = 0

	batchEntryIdx := max(0, firstDayEntryIdx+1)
	for i := batchEntryIdx; i < int32(HistoryLen); i++ {
		// Fill the history between the day of the last user activity and today
		h.Summary.History[i] = h.getCurrentEntry()
		h.updateStatistics()
	}
}
