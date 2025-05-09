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

func (h *Habits) Create(name string, stepsCount int8, stepTime int16) error {
	if len(name) > int(MaxHabitNameLength) {
		return fmt.Errorf("max habit name length cannot exceed %d", MaxHabitNameLength)
	}

	err := validateStepData(stepsCount, stepTime)

	if err != nil {
		return err
	}

	habit := newHabit(name, stepsCount, stepTime)
	h.Habits = append(h.Habits, habit)

	return nil
}

func (h *Habits) Delete(idx int) error {
	if idx < 0 || idx >= len(h.Habits) {
		return errors.New("invalid index")
	}

	h.Habits = slices.Delete(h.Habits, idx, idx+1)

	return nil
}

func (h *Habits) Print() {
	t := table.NewWriter()

	t.SetTitle("Habits")
	t.SetStyle(table.StyleLight)

	t.AppendHeader(table.Row{"#", "Name", "Checked Steps", "Steps Count", "Step Time", "Longest Streak (D)", "Total Time"})

	for idx, item := range h.Habits {
		t.AppendRow(table.Row{idx, item.Name, text.AlignCenter.Apply(getCheckedStepsText(item.CheckedSteps, item.StepsCount), 12), item.StepsCount, item.StepTime, item.LongestStreak, getTotalTimeText(item.TotalTime)})
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

func getTotalTimeText(totalTime int32) string {
	days := totalTime / 24
	hours := totalTime % 16

	text := ""

	if days > 0 {
		if days == 1 {
			text += "1 Day "
		} else {
			text += fmt.Sprintf("%d Days ", days)
		}
	}

	text += fmt.Sprintf("%d Hours", hours)

	return text
}
