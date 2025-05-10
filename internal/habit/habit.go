package habit

import (
	"fmt"
	"time"
)

const MaxHabitNameLength int8 = 16
const MaxHabitTotalTime int16 = 16 * 60 // minutes

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
	History       [6]Entry // History of last 6 days
}

type Habit struct {
	Name         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
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
		UpdatedAt:    time.Now(),
		StepsCount:   stepsCount,
		StepMinutes:  stepTime,
		CheckedSteps: 0,
		isFrozen:     false,
		Summary: Summary{
			TotalTime:     TotalTime{},
			LongestStreak: 0,
			CurrentStreak: 0,
			History:       [6]Entry{},
		},
	}

	return habit
}

func (h *Habit) CheckStep() {
	h.CheckedSteps += 1
}

func (h *Habit) UncheckStep() {
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

func (h *Habit) ChangeStepsCount(stepsCount int8) error {
	err := validateStepData(stepsCount, h.StepMinutes)

	if err != nil {
		return err
	}

	h.StepsCount = stepsCount

	return nil
}

func (h *Habit) ChangeStepMinutes(stepMinutes int16) error {
	err := validateStepData(h.StepsCount, stepMinutes)

	if err != nil {
		return err
	}

	h.StepMinutes = stepMinutes

	return nil
}

func (h *Habit) ProcessDay() {
	var entry Entry

	if h.isFrozen {
		entry = FrozenEntry{}
	} else {
		entry = ActiveEntry{Done: h.CheckedSteps, Todo: h.StepsCount}
	}

	historyLen := len(h.Summary.History)

	for i := range historyLen - 1 {
		h.Summary.History[i] = h.Summary.History[i+1]
	}

	h.Summary.History[historyLen-1] = entry
	h.Summary.TotalTime.Add(int16(h.CheckedSteps) * h.StepMinutes)

	if !h.isFrozen {
		if h.CheckedSteps >= h.StepsCount {
			h.Summary.CurrentStreak += 1

			if h.Summary.CurrentStreak > h.Summary.LongestStreak {
				h.Summary.LongestStreak = h.Summary.CurrentStreak
			}
		} else {
			h.Summary.CurrentStreak = 0
		}
	}

	h.CheckedSteps = 0
}
