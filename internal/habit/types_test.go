package habit

import (
	"fmt"
	"strings"
	"testing"
)

func TestCreate(t *testing.T) {

	t.Run("Creates a habit", func(t *testing.T) {
		habits := NewHabits()
		res := habits.Create("Test")

		if len(habits.Habits) != 1 {
			t.Error("Expected habits length to be 1, got 0")
		}

		if res != nil {
			t.Error("Expected nil, got error")
		}
	})

	t.Run(fmt.Sprintf("Returns an error when a name is longer than %d", MaxHabitNameLength), func(t *testing.T) {
		habits := NewHabits()
		res := habits.Create(strings.Repeat("A", MaxHabitNameLength+1))

		if res == nil {
			t.Error("Expected an error")
		}
	})
}

func TestDelete(t *testing.T) {

	t.Run("Deletes a habit by index", func(t *testing.T) {
		habits := NewHabits()
		habits.Create("Test")
		res := habits.Delete(0)

		if len(habits.Habits) != 0 {
			t.Errorf("Expected habits lenth to be 0, got %d", len(habits.Habits))
		}

		if res != nil {
			t.Error("Expected nil, got error")
		}
	})

	t.Run("Returns an error when the index is out of range", func(t *testing.T) {
		habits := NewHabits()
		habits.Create("Test")
		res := habits.Delete(1)

		if res == nil {
			t.Error("Expected an error")
		}
	})
}
