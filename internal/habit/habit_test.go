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

		habit.UncheckStep()
		habit.UncheckStep()

		if habit.CheckedSteps != 0 {
			t.Errorf("Expected CheckedSteps to be %d, got %d", 0, habit.CheckedSteps)
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

func TestStepMinutes(t *testing.T) {

	t.Run(fmt.Sprintf("Returns an error when the total habit time is longer than %d", MaxHabitTotalTime), func(t *testing.T) {
		habit := newHabit("Test", 1, 60)
		res := habit.ChangeStepMinutes(MaxHabitTotalTime + 1)

		if res == nil {
			t.Error("Expected an error")
		}
	})

	t.Run("Changes number of steps", func(t *testing.T) {
		habit := newHabit("Test", 1, 60)
		habit.ChangeStepMinutes(30)

		if habit.StepMinutes != 30 {
			t.Errorf("Expected StepMinutes to be %d, got %d", 30, habit.StepMinutes)
		}
	})
}
