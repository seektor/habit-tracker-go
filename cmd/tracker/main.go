package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/seektor/habit-tracker-go/internal/habit"
	"github.com/seektor/habit-tracker-go/internal/utils"
)

func main() {
	fmt.Println(utils.FgColors.FgYellow + utils.FgColors.FgBold + "=== Habit Tracker ===" + utils.FgColors.FgReset)
	fmt.Println()
	fmt.Println(getCommandText("p  [index?]", "Print all habits / a habit"))
	fmt.Println(getCommandText("a  [name] [stepTime] [stepsCount]", "Add a habit"))
	fmt.Println(getCommandText("d  [index]", "Delete a habit"))
	fmt.Println(getCommandText("sc [index] [stepTime]", "Change step time of a habit"))
	fmt.Println(getCommandText("sc [index] [stepsCount]", "Change number of steps"))
	fmt.Println(getCommandText("f  [index?]", "Freeze all habits / a habit"))
	fmt.Println(getCommandText("uf [index?]", "Unfreeze all habits / a habi"))
	fmt.Println(getCommandText("q", "Quit"))
	fmt.Println()

	fmt.Print("Enter command: ")
	reader := bufio.NewReader(os.Stdin)
	// ReadString will block until the delimiter is entered
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("An error occured while reading input. Please try again", err)
		return
	}

	// remove the delimeter from the string
	input = strings.TrimSuffix(input, "\n")
	commands := strings.Fields(input)
	if commands[0] == "q" {
		os.Exit(0)
	}
	fmt.Println(input)

	habits := habit.NewHabits()

	habits.Create("Test 1", 1, 60)
}

func getCommandText(command string, description string) string {
	return fmt.Sprintf("%s%s%s   %s", utils.FgColors.FgBold, command, utils.FgColors.FgReset, description)
}
