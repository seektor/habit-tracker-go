package habits

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/seektor/habits-tracker-go/internal/command"
	"github.com/seektor/habits-tracker-go/internal/utils"
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

func (h *Habits) Get(idx int) (*Habit, error) {
	if idx >= 0 && idx < len(h.Habits) {
		return &h.Habits[idx], nil
	}

	return nil, errors.New("invalid index")
}

func (h *Habits) Delete(idx int) error {
	if idx < 0 || idx >= len(h.Habits) {
		return errors.New("invalid index")
	}

	h.Habits = slices.Delete(h.Habits, idx, idx+1)

	return nil
}

func (h *Habits) Freeze() {
	for _, item := range h.Habits {
		item.Freeze()
	}
}

func (h *Habits) Unfreeze() {
	for _, item := range h.Habits {
		item.Unfreeze()
	}
}

func (h *Habits) UpdateToPresent() bool {
	now := time.Now()
	daysDiff := utils.GetDaysDiff(h.UpdatedAt, now)

	switch {
	case daysDiff < 0:
		utils.PrintlnError("Unknown error has occurred")
		return false
	case daysDiff == 0:
		utils.PrintlnInfo("Nothing to update")
		return false
	case daysDiff == 1:
		utils.PrintlnInfo("Updating: 1 day has passed")
	default:
		utils.PrintlnInfo(fmt.Sprintf("Updating: %d days have passed \n", daysDiff))
	}

	for idx := range h.Habits {
		h.Habits[idx].UpdateToPresent(daysDiff)
	}

	h.UpdatedAt = now
	return true
}

func (h *Habits) Print(idx int) {
	t := table.NewWriter()

	t.SetStyle(table.StyleLight)
	t.Style().Options.SeparateRows = true

	t.AppendHeader(table.Row{"#", "Name", "Checked Steps", "S Count", "S Time (min)", "Curr Streak (D)", "Lon Streak (D)", "Total Time", "History"})

	habits := h.Habits

	isSingle := idx >= 0 && idx < len(h.Habits)
	if isSingle {
		habits = habits[idx:1]
	}

	for iIdx, item := range habits {
		if isSingle {
			iIdx = idx
		}

		totalTime := item.Summary.TotalTime
		totalTime.Add(item.StepMinutes * int16(item.CheckedSteps))

		t.AppendRow(table.Row{iIdx,
			item.Name,
			text.AlignCenter.Apply(stringifyCheckedSteps(&item), 12),
			text.AlignCenter.Apply(strconv.Itoa(int(item.StepsCount)), 6),
			text.AlignCenter.Apply(strconv.Itoa(int(item.StepMinutes)), 12),
			text.AlignCenter.Apply(strconv.Itoa(int(item.Summary.CurrentStreak)), 12),
			text.AlignCenter.Apply(strconv.Itoa(int(item.Summary.LongestStreak)), 12),
			text.AlignCenter.Apply(totalTime.Stringify(), 12),
			stringifyHistory(&item),
		})
	}

	fmt.Println(t.Render())
}

func (h *Habits) PrintAll() {
	h.Print(-1)
}

func stringifyCheckedSteps(h *Habit) string {
	switch {
	case h.IsFrozen:
		return text.BgBlue.Sprint("FROZEN")
	case h.CheckedSteps < h.StepsCount:
		return text.FgRed.Sprintf("%d ❌", h.CheckedSteps)
	case h.CheckedSteps == h.StepsCount:
		return text.FgGreen.Sprintf("%d ✅", h.CheckedSteps)
	default:
		return text.FgYellow.Sprintf("%d 😎", h.CheckedSteps)
	}
}

func stringifyHistory(h *Habit) string {
	emptyBlock := "▁"
	halfBlock := "▄"
	fullBlock := "█"

	var sb strings.Builder
	history := append(h.Summary.History[:], h.getCurrentEntry())

	for idx, entry := range history {
		if entry.IsFrozen {
			sb.WriteString(utils.ColorString(utils.FgColors.Blue, halfBlock))
		} else {
			if entry.CheckedSteps == 0 {
				sb.WriteString(emptyBlock)
			} else if entry.CheckedSteps < entry.StepsCount {
				sb.WriteString(halfBlock)
			} else if entry.CheckedSteps == entry.StepsCount {
				sb.WriteString(utils.ColorString(utils.FgColors.Green, fullBlock))
			} else {
				sb.WriteString(utils.ColorString(utils.FgColors.Yellow, fullBlock))
			}
		}

		if idx != len(history)-1 {
			sb.WriteRune(' ')
		}
	}

	return sb.String()
}

var commands = []struct {
	command string
	args    string
	desc    string
}{{"p", "[index?]", "Print all habits / a habit"},
	{"a", "[name] [stepsCount] [stepMinutes]", "Add a habit"},
	{"d", "[index]", "Delete a habit"},
	{"ct", "[index] [stepMinutes]", "Change step time in minutes of a habit"},
	{"cs", "[index] [stepsCount]", "Change number of steps"},
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
		idxStr, idxStrErr := command.GetArg(0)

		if idxStrErr != nil {
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
		stepMinutesStr, stepMinutesStrErr := command.GetArg(1)

		if nameErr != nil || stepsCountStrErr != nil || stepMinutesStrErr != nil {
			utils.PrintlnError("missing arguments")
			return
		}

		stepsCount, stepsCountErr := strconv.Atoi(stepsCountStr)
		stepMinutes, stepMinutesErr := strconv.Atoi(stepMinutesStr)

		if stepsCountErr != nil || stepMinutesErr != nil {
			utils.PrintlnError("stepsCount and stepMinutes have to be a number within a proper range")
			return
		}

		err := h.Create(name, int8(stepMinutes), int16(stepsCount))
		if err != nil {
			utils.PrintlnError(err.Error())
			return
		}

		utils.PrintlnSuccess("Habit has been created")
		h.Save(utils.FileName)

	case "d":
		idxStr, idxStrErr := command.GetArg(0)

		if idxStrErr != nil {
			utils.PrintlnError("missing argument")
			return
		}

		idx, idxErr := strconv.Atoi(idxStr)

		if idxErr != nil {
			utils.PrintlnError("invalid index")
			return
		}

		err := h.Delete(idx)

		if err == nil {
			utils.PrintlnSuccess("Habit has been deleted")
			h.Save(utils.FileName)
		} else {
			utils.PrintlnError(err.Error())
		}

	case "ct":
		idxStr, idxStrErr := command.GetArg(0)
		stepMinutesStr, stepMinutesStrErr := command.GetArg(1)

		if idxStrErr != nil || stepMinutesStrErr != nil {
			utils.PrintlnError("missing arguments")
			return
		}

		idx, idxErr := strconv.Atoi(idxStr)
		stepMinutes, stepMinutesErr := strconv.Atoi(stepMinutesStr)

		if idxErr != nil {
			utils.PrintlnError("invalid index")
			return
		}

		if stepMinutesErr != nil {
			utils.PrintlnError("invalid number of minutes")
			return
		}

		habit, habitErr := h.Get(idx)

		if habitErr != nil {
			utils.PrintlnError(habitErr.Error())
			return
		}

		err := habit.SetStepMinutes(int16(stepMinutes))

		if err == nil {
			utils.PrintlnSuccess("Step time has been updated")
			h.Save(utils.FileName)
		} else {
			utils.PrintlnError(err.Error())
		}

	case "cs":
		idxStr, idxStrErr := command.GetArg(0)
		stepsCountStr, stepsCountStrErr := command.GetArg(1)

		if idxStrErr != nil || stepsCountStrErr != nil {
			utils.PrintlnError("missing arguments")
			return
		}

		idx, idxErr := strconv.Atoi(idxStr)
		stepsCount, stepsCountErr := strconv.Atoi(stepsCountStr)

		if idxErr != nil {
			utils.PrintlnError("invalid index")
			return
		}

		if stepsCountErr != nil {
			utils.PrintlnError("invalid number of steps")
			return
		}

		habit, habitErr := h.Get(idx)

		if habitErr != nil {
			utils.PrintlnError(habitErr.Error())
			return
		}

		err := habit.SetStepsCount(int8(stepsCount))

		if err == nil {
			utils.PrintlnSuccess("Steps count has been updated")
			h.Save(utils.FileName)
		} else {
			utils.PrintlnError(err.Error())
		}

	case "f":
		idxStr, idxStrErr := command.GetArg(0)

		if idxStrErr != nil {
			h.Freeze()
			utils.PrintlnSuccess("Habits have been frozen")
			h.Save(utils.FileName)
		} else {
			idx, err := strconv.Atoi(idxStr)
			if err != nil {
				utils.PrintlnError("invalid index")
				return
			}

			habit, habitErr := h.Get(idx)

			if habitErr != nil {
				utils.PrintlnError(habitErr.Error())
				return
			}
			habit.Freeze()
			utils.PrintlnSuccess("Habit has been frozen")
			h.Save(utils.FileName)
		}

	case "uf":
		idxStr, idxStrErr := command.GetArg(0)

		if idxStrErr != nil {
			h.Unfreeze()
			utils.PrintlnSuccess("Habits have been unfrozen")
			h.Save(utils.FileName)
		} else {
			idx, err := strconv.Atoi(idxStr)
			if err != nil {
				utils.PrintlnError("invalid index")
				return
			}

			habit, habitErr := h.Get(idx)

			if habitErr != nil {
				utils.PrintlnError(habitErr.Error())
				return
			}
			habit.Freeze()
			utils.PrintlnSuccess("Habit has been unfrozen")
			h.Save(utils.FileName)
		}

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
