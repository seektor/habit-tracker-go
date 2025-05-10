package habit

import (
	"fmt"
	"testing"
)

func TestCheckStep(t *testing.T) {

	t.Run("checks a habit step", func(t *testing.T) {
		habit := newHabit("Test", 1, 60)
		habit.CheckStep()
		habit.CheckStep()
		habit.CheckStep()

		if habit.CheckedSteps != 3 {
			t.Errorf("expected CheckedSteps to be %d, got %d", 3, habit.CheckedSteps)
		}
	})

	t.Run("does not check a habit when it is frozen", func(t *testing.T) {
		habit := newHabit("Test", 1, 60)
		habit.Freeze()
		habit.CheckStep()

		if habit.CheckedSteps != 0 {
			t.Errorf("expected CheckedSteps to be %d, got %d", 0, habit.CheckedSteps)
		}
	})
}

func TestUncheckStep(t *testing.T) {

	t.Run("unchecks a habit step", func(t *testing.T) {
		habit := newHabit("Test", 1, 60)
		habit.CheckStep()
		habit.CheckStep()
		habit.UncheckStep()

		if habit.CheckedSteps != 1 {
			t.Errorf("expected CheckedSteps to be %d, got %d", 1, habit.CheckedSteps)
		}

		habit.UncheckStep()
		habit.UncheckStep()

		if habit.CheckedSteps != 0 {
			t.Errorf("expected CheckedSteps to be %d, got %d", 0, habit.CheckedSteps)
		}
	})
}

func TestSetStepsCount(t *testing.T) {

	t.Run(fmt.Sprintf("returns an error when the total habit time is longer than %d", MaxHabitTotalTime), func(t *testing.T) {
		habit := newHabit("Test", 1, 60)
		res := habit.SetStepsCount(17)

		if res == nil {
			t.Error("expected an error")
		}
	})

	t.Run("sets number of steps", func(t *testing.T) {
		habit := newHabit("Test", 1, 60)
		habit.SetStepsCount(3)

		if habit.StepsCount != 3 {
			t.Errorf("expected StepsCount to be %d, got %d", 3, habit.StepsCount)
		}
	})
}

func TestSetStepMinutes(t *testing.T) {

	t.Run(fmt.Sprintf("returns an error when the total habit time is longer than %d", MaxHabitTotalTime), func(t *testing.T) {
		habit := newHabit("Test", 1, 60)
		res := habit.SetStepMinutes(MaxHabitTotalTime + 1)

		if res == nil {
			t.Error("expected an error")
		}
	})

	t.Run("sets number of steps", func(t *testing.T) {
		habit := newHabit("Test", 1, 60)
		habit.SetStepMinutes(30)

		if habit.StepMinutes != 30 {
			t.Errorf("expected StepMinutes to be %d, got %d", 30, habit.StepMinutes)
		}
	})
}

func TestFreeze(t *testing.T) {
	t.Run("freezes the habit and clears the CheckedSteps count", func(t *testing.T) {
		habit := newHabit("Test", 1, 60)
		habit.CheckStep()
		habit.Freeze()

		if habit.isFrozen != true {
			t.Errorf("expected isFrozen to be %t, got %t", true, habit.isFrozen)
		}

		if habit.CheckedSteps != 0 {
			t.Errorf("expected CheckedSteps to be %d, got %d", 0, habit.CheckedSteps)
		}
	})
}

func TestUnfreeze(t *testing.T) {
	t.Run("unfreezes the habit", func(t *testing.T) {
		habit := newHabit("Test", 1, 60)
		habit.Unfreeze()

		if habit.isFrozen == true {
			t.Error("habit is not unfrozen correctly")
		}
	})
}
