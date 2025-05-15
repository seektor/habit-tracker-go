package habits

import (
	"fmt"
	"time"
)

const MaxHabitNameLength int8 = 16
const MaxHabitTotalTime int16 = 16 * 60 // minutes
const HistoryLen int8 = 6

type Entry interface {
	isEntry()
}

type ActiveEntry struct {
	Done int8
	Todo int8
}

func (ae ActiveEntry) isEntry() {}

type FrozenEntry struct{}

func (fe FrozenEntry) isEntry() {}

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
	isFrozen     bool
	Summary      Summary
}

func newHabit(name string, stepsCount int8, stepTime int16) Habit {
	habit := Habit{
		Name:         name,
		CreatedAt:    time.Now(),
		StepsCount:   stepsCount,
		StepMinutes:  stepTime,
		CheckedSteps: 0,
		isFrozen:     false,
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
	if h.isFrozen {
		return
	}

	h.CheckedSteps += 1
}

func (h *Habit) UncheckStep() {
	if h.isFrozen {
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
	h.isFrozen = true
	h.CheckedSteps = 0
}

func (h *Habit) Unfreeze() {
	h.isFrozen = false
}

func (h *Habit) getCurrentEntry() Entry {
	if h.isFrozen {
		return FrozenEntry{}
	} else {
		return ActiveEntry{
			Done: h.CheckedSteps,
			Todo: h.StepsCount,
		}
	}
}

func (h *Habit) updateStatistics() {
	if h.isFrozen {
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
