package habit

import (
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
	StepTime  int8 // Minutes
	Entries   [7]Entry
	Summary   *Summary
}

func newHabit(name string) *Habit {
	habit := &Habit{
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		StepCount: 1,
		StepTime:  60,
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

func (h *Habits) Add(name string) {
	// TODO: Validate len
	habit := newHabit(name)
	h.Habits = append(h.Habits, habit)
}
