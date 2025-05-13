package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/seektor/habit-tracker-go/internal/habit"
	"github.com/seektor/habit-tracker-go/internal/utils"
)

func main() {
	fmt.Println(utils.FgColors.Yellow +
		utils.FgColors.Bold +
		"=== Habit Tracker ===" +
		utils.FgColors.Reset)

	habits := habit.NewHabits()

	habits.Create("Test 1", 1, 60)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println()
		fmt.Print("Enter command: ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occured while reading input. Please try again", err)
			return
		}

		command, args := parseInput(input)

		handleCommand(&habits, command, args)
	}
}

func parseInput(input string) (string, []string) {
	inputArgs := strings.Fields(input)
	command := ""
	args := []string{}

	if len(inputArgs) >= 1 {
		command = inputArgs[0]
	}

	if len(inputArgs) > 1 {
		args = inputArgs[1:]
	}

	return command, args
}

func printCommands() {
	fmt.Println()
	fmt.Println(getCommandText("p  [index?]", "Print all habits / a habit"))
	fmt.Println(getCommandText("a  [name] [stepTime] [stepsCount]", "Add a habit"))
	fmt.Println(getCommandText("d  [index]", "Delete a habit"))
	fmt.Println(getCommandText("ct [index] [stepTime]", "Change step time of a habit"))
	fmt.Println(getCommandText("cs [index] [stepsCount]", "Change number of steps"))
	fmt.Println(getCommandText("f  [index?]", "Freeze all habits / a habit"))
	fmt.Println(getCommandText("uf [index?]", "Unfreeze all habits / a habit"))
	fmt.Println(getCommandText("q", "Quit"))
	fmt.Println()
}

func handleCommand(habits *habit.Habits, command string, args []string) {
	switch command {
	case "p":
		idx, err := getAtIdx(args, 0)
		if err != nil {
			habits.Print()
		} else {
			fmt.Print(idx)
		}
	case "q":
		os.Exit(0)
	default:
		fmt.Println(utils.FgColors.Red + "Unknown command" + utils.FgColors.Reset)
		printCommands()
	}

}

func getAtIdx(slice []string, idx int) (string, error) {
	if idx >= 0 && idx < len(slice) {
		return slice[idx], nil
	}

	return "", errors.New("no value")
}

func getCommandText(command string, description string) string {
	return fmt.Sprintf("%s%s%s   %s", utils.FgColors.Bold, command, utils.FgColors.Reset, description)
}
