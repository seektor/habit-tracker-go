package habit

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/seektor/habit-tracker-go/internal/utils"
)

type Habits struct {
	Habits    []Habit
	UpdatedAt time.Time
}

func NewHabits() Habits {
	return Habits{
		Habits:    make([]Habit, 0),
		UpdatedAt: time.Now(),
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

func (h *Habits) Freeze(idx int) {
	for _, item := range h.Habits {
		item.Freeze()
	}
}

func (h *Habits) Unfreeze(idx int) {
	for _, item := range h.Habits {
		item.Unfreeze()
	}
}

func (h *Habits) UpdateToPresent() {
	now := time.Now()
	daysDiff := utils.GetDaysDiff(h.UpdatedAt, now)

	switch {
	case daysDiff < 0:
		fmt.Println("Unknown error has occurred")
		return
	case daysDiff == 0:
		fmt.Println("=== Nothing to update ===")
		return
	case daysDiff == 1:
		fmt.Println("=== Updating: 1 day has passed ===")
	default:
		fmt.Printf("=== Updating: %d days have passed ===\n", daysDiff)
	}

	for idx := range h.Habits {
		h.Habits[idx].UpdateToPresent(daysDiff)
	}

	h.UpdatedAt = now
}

func (h *Habits) Print() {
	t := table.NewWriter()

	t.SetTitle("Habits")
	t.SetStyle(table.StyleLight)

	t.AppendHeader(table.Row{"#", "Name", "Checked Steps", "Steps Count", "Step Time (min)", "Longest Streak (D)", "Total Time"})

	for idx, item := range h.Habits {
		totalTime := item.Summary.TotalTime
		totalTime.Add(item.StepMinutes * int16(item.CheckedSteps))

		t.AppendRow(table.Row{idx,
			item.Name,
			text.AlignCenter.Apply(getCheckedStepsText(&item), 12),
			item.StepsCount,
			item.StepMinutes,
			item.Summary.LongestStreak,
			totalTime.Stringify(),
		})
	}

	fmt.Println(t.Render())
}

func getCheckedStepsText(h *Habit) string {
	switch {
	case h.isFrozen:
		return text.BgBlue.Sprint("FROZEN")
	case h.CheckedSteps < h.StepsCount:
		return text.FgRed.Sprintf("%d ❌", h.CheckedSteps)
	case h.CheckedSteps == h.StepsCount:
		return text.FgGreen.Sprintf("%d ✅", h.CheckedSteps)
	default:
		return text.FgBlue.Sprintf("%d 😎", h.CheckedSteps)
	}
}
