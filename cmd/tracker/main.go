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

const FileName = "habit_tracker.json"

func main() {
	fmt.Println(utils.FgColors.Yellow + utils.FgColors.Bold +
		"=== Habit Tracker ===" +
		utils.FgColors.Reset)

	habits := habit.NewHabits()

	if err := habits.Load(FileName); err != nil {
		fmt.Println(utils.FgColors.Red + err.Error())
		os.Exit(1)
	}

	fmt.Println()

	if len(habits.Habits) > 0 {
		habits.PrintAll()
	} else {
		habits.PrintCommands()
	}

	fmt.Println()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter command: ")

		input, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println("An error occured while reading input. Please try again", err)
			return
		}

		command, args := parseInput(input)

		handleCommand(habits, command, args)
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

func handleCommand(habits *habit.Habits, command string, args []string) {
	switch command {
	case "p":
		idx, err := getAtIdx(args, 0)
		if err != nil {
			habits.PrintAll()
		} else {
			fmt.Print(idx)
		}
	case "q":
		os.Exit(0)
	default:
		fmt.Println(utils.FgColors.Red + "=== Unknown command ===" + utils.FgColors.Reset)
		habits.PrintCommands()
	}

}

func getAtIdx(slice []string, idx int) (string, error) {
	if idx >= 0 && idx < len(slice) {
		return slice[idx], nil
	}

	return "", errors.New("no value")
}
