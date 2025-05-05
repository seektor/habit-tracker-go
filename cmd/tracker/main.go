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
	// Implement Delete
	printHabits(habits)
}

func printHabits(habits habit.Habits) {
	fmt.Println("-- Printing Habits --")

	res, err := json.MarshalIndent(habits, "", "  ")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(res))
}
