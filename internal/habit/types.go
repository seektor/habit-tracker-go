package habit

import (
	"errors"
	"fmt"
	"slices"
	"time"
)

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
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	StepCount int8
	StepTime  int16 // Minutes
	Entries   [7]Entry
	Summary   *Summary
}

func newHabit(name string, stepCount int8, stepTime int16) *Habit {
	habit := &Habit{
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		StepCount: stepCount,
		StepTime:  stepTime,
		Entries:   [7]Entry{},
		Summary: &Summary{
			LongestStreak: 0,
			TotalTime:     0,
		},
	}

	return habit
}

type Habits struct {
	Habits []*Habit
}

func NewHabits() Habits {
	return Habits{
		Habits: make([]*Habit, 0),
	}
}

const MaxHabitNameLength int8 = 16
const MaxHabitTotalTime int16 = 60 * 16

func (h *Habits) Create(name string, stepCount int8, stepTime int16) error {
	if len(name) > int(MaxHabitNameLength) {
		return fmt.Errorf("Max habit name length cannot exceed %d", MaxHabitNameLength)
	}

	if stepCount < 1 || stepTime < 1 {
		return fmt.Errorf("Step count and Step time has to be a positive value")
	}

	if int16(stepCount)*stepTime > MaxHabitTotalTime {
		return fmt.Errorf("Max habit total time cannot exceed %d", MaxHabitTotalTime)
	}

	habit := newHabit(name, stepCount, stepTime)
	h.Habits = append(h.Habits, habit)

	return nil
}

func (h *Habits) Delete(idx int) error {
	if idx < 0 || idx >= len(h.Habits) {
		return errors.New("Invalid index")
	}

	h.Habits = slices.Delete(h.Habits, idx, idx+1)

	return nil
}
