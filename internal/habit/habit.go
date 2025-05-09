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

type Habit struct {
	Name          string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	StepsCount    int8
	StepMinutes   int16 // minutes
	CheckedSteps  int8
	TotalTime     *TotalTime
	LongestStreak int16
	Entries       [6]Entry // History of last 6 days
}

func newHabit(name string, stepsCount int8, stepTime int16) *Habit {
	habit := &Habit{
		Name:          name,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		StepsCount:    stepsCount,
		StepMinutes:   stepTime,
		CheckedSteps:  0,
		Entries:       [6]Entry{},
		LongestStreak: 0,
		TotalTime:     &TotalTime{},
	}

	return habit
}

func (h *Habit) CheckStep() {
	h.CheckedSteps += 1
	h.TotalTime.Add(h.StepMinutes)
}

func (h *Habit) UncheckStep() {
	if h.CheckedSteps > 0 {
		h.CheckedSteps -= 1
		h.TotalTime.Subtract(h.StepMinutes)
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

func (h *Habit) ChangeStepTime(stepCount int8) {
}
