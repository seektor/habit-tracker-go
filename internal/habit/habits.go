package habit

import (
	"errors"
	"fmt"
	"slices"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type Habits struct {
	Habits []*Habit
}

func NewHabits() Habits {
	return Habits{
		Habits: make([]*Habit, 0),
	}
}

func (h *Habits) Create(name string, stepCount int8, stepTime int8) error {
	if len(name) > int(MaxHabitNameLength) {
		return fmt.Errorf("Max habit name length cannot exceed %d", MaxHabitNameLength)
	}

	err := validateStepData(stepCount, stepTime)

	if err != nil {
		return err
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

func (h *Habits) Print() {
	t := table.NewWriter()

	t.SetTitle("Habits")
	t.SetStyle(table.StyleLight)

	t.AppendHeader(table.Row{"#", "Name", "Checked Steps", "Steps", "Step Time", "Longest Streak (D)", "Total Time (h)"})

	for idx, item := range h.Habits {
		t.AppendRow(table.Row{idx, item.Name, text.AlignCenter.Apply(getCheckedStepsText(item.CheckedSteps, item.Steps), 12), item.Steps, item.StepTime, item.Summary.LongestStreak, item.Summary.TotalTime})
	}

	fmt.Println(t.Render())
}

func getCheckedStepsText(checkedSteps int8, steps int8) string {
	switch {
	case checkedSteps < steps:
		return text.FgRed.Sprintf("%d ❌", checkedSteps)
	case checkedSteps == steps:
		return text.FgGreen.Sprintf("%d ✅", checkedSteps)
	default:
		return text.FgBlue.Sprintf("%d 😎", checkedSteps)
	}
}
