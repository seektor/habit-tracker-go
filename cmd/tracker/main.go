package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/seektor/habit-tracker-go/internal/habit"
)

func main() {
	fmt.Println("Hello World")

	habits := habit.NewHabits()

	habits.Create("Test 1", 1, 60)

	habits.Habits[0].CheckStep()
	habits.Habits[0].Freeze()
	habits.Print()
	// habits.Habits[0].ProcessDay()
	fmt.Printf("%v", habits.Habits[0].Summary.History)
}

func printHabits(habits habit.Habits) {
	fmt.Println("-- Printing Habits --")

	res, err := json.MarshalIndent(habits, "", "  ")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(res))
}
