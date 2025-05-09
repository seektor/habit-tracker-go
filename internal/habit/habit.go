package habit

import (
	"fmt"
	"time"
)

const MaxHabitNameLength int8 = 16
const MaxHabitTotalTime int8 = 16

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
	LongestStreak int16
	TotalTime     int32
}

type Habit struct {
	Name         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Steps        int8
	StepTime     int8
	CheckedSteps int8
	Entries      [7]Entry
	Summary      *Summary
}

func newHabit(name string, stepsCount int8, stepTime int8) *Habit {
	habit := &Habit{
		Name:         name,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Steps:        stepsCount,
		StepTime:     stepTime,
		CheckedSteps: 0,
		Entries:      [7]Entry{},
		Summary: &Summary{
			LongestStreak: 0,
			TotalTime:     0,
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

func validateStepData(stepCount int8, stepTime int8) error {
	if stepCount < 1 || stepTime < 1 {
		return fmt.Errorf("Step count and Step time has to be a positive value")
	}

	if stepCount*stepTime > MaxHabitTotalTime {
		return fmt.Errorf("Max habit total time cannot exceed %d", MaxHabitTotalTime)
	}

	return nil
}

func (h *Habit) ChangeStepCount(stepCount int8) {
}
