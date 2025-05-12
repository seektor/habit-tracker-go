package habit

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {

	t.Run("creates a habit", func(t *testing.T) {
		habits := NewHabits()
		res := habits.Create("Test", 1, 60)

		if len(habits.Habits) != 1 {
			t.Error("expected habits length to be 1, got 0")
		}

		if res != nil {
			t.Error("expected nil, got error")
		}
	})

	t.Run(fmt.Sprintf("returns an error when a name is longer than %d", MaxHabitNameLength), func(t *testing.T) {
		habits := NewHabits()
		res := habits.Create(strings.Repeat("A", int(MaxHabitNameLength)+1), 1, 60)

		if res == nil {
			t.Error("expected an error")
		}
	})

	t.Run("returns an error when the StepCount or StepTime are smaller than 1", func(t *testing.T) {
		habits := NewHabits()
		res := habits.Create("Test", -1, 0)

		if res == nil {
			t.Error("expected an error")
		}
	})

	t.Run(fmt.Sprintf("returns an error when the total habit time is longer than %d", MaxHabitTotalTime), func(t *testing.T) {
		habits := NewHabits()
		res := habits.Create("Test", 17, 60)

		if res == nil {
			t.Error("expected an error")
		}
	})
}

func TestDelete(t *testing.T) {

	t.Run("deletes a habit by index", func(t *testing.T) {
		habits := NewHabits()
		habits.Create("Test", 1, 60)
		res := habits.Delete(0)

		if len(habits.Habits) != 0 {
			t.Errorf("expected habits length to be 0, got %d", len(habits.Habits))
		}

		if res != nil {
			t.Error("expected nil, got error")
		}
	})

	t.Run("returns an error when the index is out of range", func(t *testing.T) {
		habits := NewHabits()
		habits.Create("Test", 1, 60)
		res := habits.Delete(1)

		if res == nil {
			t.Error("expected an error")
		}
	})
}

func TestHabitsUpdateToPresent(t *testing.T) {

	t.Run("does not update habits when there is no day difference", func(t *testing.T) {
		habits := NewHabits()
		habits.Create("Test", 1, 60)

		preUpdatedHabit := habits.Habits[0]

		habits.UpdateToPresent()

		if habits.Habits[0] != preUpdatedHabit {
			t.Error("habits have not been updated")
		}

		now := time.Now()

		isDateUpdated := habits.UpdatedAt.Year() == now.Year() &&
			habits.UpdatedAt.Month() == now.Month() &&
			habits.UpdatedAt.Day() == now.Day()

		if !isDateUpdated {
			t.Error("date have not been updated")
		}
	})

	t.Run("updates habits when there is a day difference", func(t *testing.T) {
		habits := NewHabits()
		habits.Create("Test", 1, 60)
		habits.UpdatedAt = habits.UpdatedAt.AddDate(0, 0, -1)

		preUpdatedHabit := habits.Habits[0]

		habits.UpdateToPresent()
		fmt.Printf("%v %v", preUpdatedHabit.Summary.History, habits.Habits[0].Summary.History)

		if habits.Habits[0] == preUpdatedHabit {
			t.Error("habits have not been updated")
		}

		now := time.Now()

		isDateUpdated := habits.UpdatedAt.Year() == now.Year() &&
			habits.UpdatedAt.Month() == now.Month() &&
			habits.UpdatedAt.Day() == now.Day()

		if !isDateUpdated {
			t.Error("date have not been updated")
		}
	})
}
