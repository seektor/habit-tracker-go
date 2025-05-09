package habit

import (
	"fmt"
	"testing"
)

func TestCheckStep(t *testing.T) {

	t.Run("Checks a habit step", func(t *testing.T) {
		habit := newHabit("Test", 1, 60)
		habit.CheckStep()
		habit.CheckStep()
		habit.CheckStep()

		if habit.CheckedSteps != 3 {
			t.Errorf("Expected CheckedSteps to be %d, got %d", 3, habit.CheckedSteps)
		}

		want := TotalTime{Hours: 3}
		if *habit.TotalTime != want {
			t.Errorf("Expected TotalTime to be %v, got %d", want, *habit.TotalTime)
		}
	})
}

func TestUncheckStep(t *testing.T) {

	t.Run("Unchecks a habit step", func(t *testing.T) {
		habit := newHabit("Test", 1, 60)
		habit.CheckStep()
		habit.CheckStep()
		habit.UncheckStep()

		if habit.CheckedSteps != 1 {
			t.Errorf("Expected CheckedSteps to be %d, got %d", 1, habit.CheckedSteps)
		}

		want := TotalTime{Hours: 1}
		if *habit.TotalTime != want {
			t.Errorf("Expected TotalTime to be %d, got %d", want, *habit.TotalTime)
		}

		habit.UncheckStep()
		habit.UncheckStep()

		if habit.CheckedSteps != 0 {
			t.Errorf("Expected CheckedSteps to be %d, got %d", 0, habit.CheckedSteps)
		}

		want = TotalTime{}
		if *habit.TotalTime != want {
			t.Errorf("Expected TotalTime to be %d, got %d", want, *habit.TotalTime)
		}
	})
}

func TestChangeStepsCount(t *testing.T) {

	t.Run(fmt.Sprintf("Returns an error when the total habit time is longer than %d", MaxHabitTotalTime), func(t *testing.T) {
		habit := newHabit("Test", 1, 60)
		res := habit.ChangeStepsCount(17)

		if res == nil {
			t.Error("Expected an error")
		}
	})

	t.Run("Changes number of steps", func(t *testing.T) {
		habit := newHabit("Test", 1, 60)
		habit.ChangeStepsCount(3)

		if habit.StepsCount != 3 {
			t.Errorf("Expected StepsCount to be %d, got %d", 3, habit.StepsCount)
		}
	})
}
