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

	habits.Add("Test 1")
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
