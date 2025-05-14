package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/seektor/habit-tracker-go/internal/command"
	"github.com/seektor/habit-tracker-go/internal/habit"
	"github.com/seektor/habit-tracker-go/internal/utils"
)

func main() {
	fmt.Println(utils.FgColors.Yellow + utils.FgColors.Bold +
		"=== Habit Tracker ===" +
		utils.FgColors.Reset)

	habits := habit.NewHabits()

	if err := habits.Load(); err != nil {
		fmt.Println(utils.FgColors.Red + err.Error())
		os.Exit(1)
	}

	fmt.Println()

	if len(habits.Habits) > 0 {
		habits.PrintAll()
	} else {
		habits.PrintCommands()
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println()
		fmt.Print("Enter command: ")

		input, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println("An error occured while reading input. Please try again", err)
			return
		}

		command := command.NewCommand(input)
		habits.Execute(command)
	}
}
