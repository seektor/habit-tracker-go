package habit

import (
	"fmt"
	"strings"
	"testing"
)

func TestCreate(t *testing.T) {

	t.Run("Creates a habit", func(t *testing.T) {
		habits := NewHabits()
		res := habits.Create("Test", 1, 1)

		if len(habits.Habits) != 1 {
			t.Error("Expected habits length to be 1, got 0")
		}

		if res != nil {
			t.Error("Expected nil, got error")
		}
	})

	t.Run(fmt.Sprintf("Returns an error when a name is longer than %d", MaxHabitNameLength), func(t *testing.T) {
		habits := NewHabits()
		res := habits.Create(strings.Repeat("A", int(MaxHabitNameLength)+1), 1, 1)

		if res == nil {
			t.Error("Expected an error")
		}
	})

	t.Run("Returns an error when the StepCount or StepTime are smaller than 1", func(t *testing.T) {
		habits := NewHabits()
		res := habits.Create("Test", -1, 0)

		if res == nil {
			t.Error("Expected an error")
		}
	})

	t.Run(fmt.Sprintf("Returns an error when the total habit time is longer than %d", MaxHabitTotalTime), func(t *testing.T) {
		habits := NewHabits()
		res := habits.Create("Test", -1, 0)

		if res == nil {
			t.Error("Expected an error")
		}
	})
}

func TestDelete(t *testing.T) {

	t.Run("Deletes a habit by index", func(t *testing.T) {
		habits := NewHabits()
		habits.Create("Test", 1, 1)
		res := habits.Delete(0)

		if len(habits.Habits) != 0 {
			t.Errorf("Expected habits length to be 0, got %d", len(habits.Habits))
		}

		if res != nil {
			t.Error("Expected nil, got error")
		}
	})

	t.Run("Returns an error when the index is out of range", func(t *testing.T) {
		habits := NewHabits()
		habits.Create("Test", 1, 1)
		res := habits.Delete(1)

		if res == nil {
			t.Error("Expected an error")
		}
	})
}

func TestCheckStep(t *testing.T) {

	t.Run("Checks a habit step", func(t *testing.T) {
		habit := newHabit("Test", 1, 1)
		habit.CheckStep()
		habit.CheckStep()
		habit.CheckStep()

		if habit.CheckedSteps != 3 {
			t.Errorf("Expected CheckedSteps to be %d, got %d", 3, habit.CheckedSteps)
		}
	})
}

func TestUncheckStep(t *testing.T) {

	t.Run("Unchecks a habit step", func(t *testing.T) {
		habit := newHabit("Test", 1, 1)
		habit.CheckStep()
		habit.CheckStep()
		habit.UncheckStep()

		if habit.CheckedSteps != 1 {
			t.Errorf("Expected CheckedSteps to be %d, got %d", 1, habit.CheckedSteps)
		}

		habit.UncheckStep()
		habit.UncheckStep()

		if habit.CheckedSteps != 0 {
			t.Errorf("Expected CheckedSteps to be %d, got %d", 0, habit.CheckedSteps)
		}
	})
}
