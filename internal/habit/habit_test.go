package habit

import (
	"testing"
)

func TestCheckStep(t *testing.T) {

	t.Run("Checks a habit step", func(t *testing.T) {
		habit := newHabit("Test", 1, 1)
		habit.CheckStep()
		habit.CheckStep()
		habit.CheckStep()

		if habit.CheckedSteps != 3 {
			t.Errorf("Expected CheckedSteps to be %d, got %d", 3, habit.CheckedSteps)
		}

		if habit.TotalTime != 3 {
			t.Errorf("Expected TotalTime to be %d, got %d", 3, habit.TotalTime)
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

		if habit.TotalTime != 1 {
			t.Errorf("Expected TotalTime to be %d, got %d", 1, habit.TotalTime)
		}

		habit.UncheckStep()
		habit.UncheckStep()

		if habit.CheckedSteps != 0 {
			t.Errorf("Expected CheckedSteps to be %d, got %d", 0, habit.CheckedSteps)
		}

		if habit.TotalTime != 0 {
			t.Errorf("Expected TotalTime to be %d, got %d", 0, habit.TotalTime)
		}
	})
}
