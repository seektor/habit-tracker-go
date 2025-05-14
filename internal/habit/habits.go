package habit

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/seektor/habit-tracker-go/internal/command"
	"github.com/seektor/habit-tracker-go/internal/utils"
)

type Habits struct {
	Habits    []Habit
	UpdatedAt time.Time
}

func NewHabits() *Habits {
	return &Habits{
		Habits:    make([]Habit, 0),
		UpdatedAt: time.Now(),
	}
}

func (h *Habits) Load() error {
	file, err := os.ReadFile(utils.FileName)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return err
	}

	if len(file) == 0 {
		return err
	}

	err = json.Unmarshal(file, h)

	if err != nil {
		return err
	}

	return nil
}

func (h *Habits) Save(filename string) error {
	data, err := json.Marshal(h)

	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
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

func (h *Habits) Print(idx int) {
	t := table.NewWriter()

	t.SetStyle(table.StyleLight)

	t.AppendHeader(table.Row{"#", "Name", "Checked Steps", "Steps Count", "Step Time (min)", "Longest Streak (D)", "Total Time"})

	habits := h.Habits

	if idx >= 0 && idx < len(h.Habits) {
		t.SetTitle("Habit %d", idx)
		habits = habits[idx:1]
	} else {
		t.SetTitle("Habits")
	}

	for idx, item := range habits {
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

func (h *Habits) PrintAll() {
	h.Print(-1)
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

var commands = []struct {
	command string
	args    string
	desc    string
}{{"p", "[index?]", "Print all habits / a habit"},
	{"a", "[name] [stepsCount] [stepTime]", "Add a habit"},
	{"d", "[index]", "Delete a habit"},
	{"ct", "[index] [stepsTime]", "Change step time of a habit"},
	{"cs", "[index] [stepCount]", "Change number of steps"},
	{"f", "[index]?", "Freeze all habits / a habit"},
	{"uf", "[index]?", "Unfreeze all habits / a habit"},
	{"q", "", "Quit"},
}

func (h *Habits) PrintCommands() {
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false

	for _, item := range commands {

		t.AppendRow(table.Row{text.Bold.Sprint(item.command),
			text.Bold.Sprint(item.args),
			item.desc,
		})
	}

	fmt.Println(t.Render())
}

func (h *Habits) Execute(command command.Command) {
	switch command.Command {
	case "p":
		idxStr, err := command.GetArg(0)

		if err != nil {
			h.PrintAll()
		} else {
			idx, err := strconv.Atoi(idxStr)
			if err != nil {
				utils.PrintlnError("invalid index")
				return
			}
			h.Print(int(idx))
		}

	case "a":
		name, nameErr := command.GetArg(0)
		stepsCountStr, stepsCountStrErr := command.GetArg(2)
		stepTimeStr, stepTimeStrErr := command.GetArg(1)

		if nameErr != nil || stepsCountStrErr != nil || stepTimeStrErr != nil {
			utils.PrintlnError("missing arguments")
			return
		}

		stepsCount, stepsCountErr := strconv.Atoi(stepsCountStr)
		stepTime, stepTimeErr := strconv.Atoi(stepTimeStr)

		if stepsCountErr != nil || stepTimeErr != nil {
			utils.PrintlnError("stepsCount and stepTime have to be a number within a proper range")
			return
		}

		err := h.Create(name, int8(stepTime), int16(stepsCount))
		if err != nil {
			utils.PrintlnError(err.Error())
			return
		}

		utils.PrintlnSuccess("Habit has been created")
		h.Save(utils.FileName)

	case "q":
		utils.PrintlnSuccess("Bye bye")
		os.Exit(0)

	default:
		fmt.Println()
		utils.PrintlnError("unknown command")
		fmt.Println()
		h.PrintCommands()
	}

}
