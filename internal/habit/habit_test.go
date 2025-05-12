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

func isUpdated(habit *Habit, currentStreak int16, longestStreak int16, totalTime TotalTime, history [HistoryLen]Entry) bool {
	isStreakUpdated := habit.Summary.CurrentStreak == currentStreak
	isLongestStreakUpdated := habit.Summary.LongestStreak == longestStreak
	isCheckedStepUpdated := habit.CheckedSteps == 0
	isTotalTimeUpdated := habit.Summary.TotalTime == totalTime
	isHistoryUpdated := habit.Summary.History == history

	return isStreakUpdated &&
		isLongestStreakUpdated &&
		isCheckedStepUpdated &&
		isTotalTimeUpdated &&
		isHistoryUpdated
}

func TestHabitUpdateToPresent(t *testing.T) {

	t.Run("does not update when there is no day change", func(t *testing.T) {
		habit := newHabit("Test", 1, 60)
		habit.CheckStep()
		habit.UpdateToPresent(0)

		isUpdated := isUpdated(&habit, 0, 0, TotalTime{Hours: 0}, [HistoryLen]Entry{nil, nil, nil, nil, nil, nil})

		if isUpdated {
			t.Error("habit has not been updated successfully")
		}
	})

	t.Run("updates habit by 1 day", func(t *testing.T) {
		habit := newHabit("Test", 2, 60)
		habit.CheckStep()
		habit.CheckStep()
		habit.UpdateToPresent(1)

		isUpdated := isUpdated(&habit, 1, 1, TotalTime{Hours: 2}, [HistoryLen]Entry{nil, nil, nil, nil, nil, ActiveEntry{2, 2}})

		if !isUpdated {
			t.Error("habit has not been updated successfully")
		}
	})

	t.Run("updates frozen habit by 1 day", func(t *testing.T) {
		habit := newHabit("Test", 2, 60)
		habit.CheckStep()
		habit.CheckStep()
		habit.Summary.CurrentStreak = 2
		habit.Summary.LongestStreak = 3
		habit.Freeze()
		habit.UpdateToPresent(1)

		isUpdated := isUpdated(&habit, 2, 3, TotalTime{}, [HistoryLen]Entry{nil, nil, nil, nil, nil, FrozenEntry{}})

		if !isUpdated {
			t.Error("habit has not been updated successfully")
		}
	})

	t.Run("updates habit by 3 days", func(t *testing.T) {
		habit := newHabit("Test", 2, 60)
		habit.CheckStep()
		habit.CheckStep()
		habit.UpdateToPresent(3)

		isUpdated := isUpdated(&habit, 0, 1, TotalTime{Hours: 2}, [HistoryLen]Entry{nil, nil, nil, ActiveEntry{2, 2}, ActiveEntry{0, 2}, ActiveEntry{0, 2}})

		if !isUpdated {
			t.Error("habit has not been updated successfully")
		}
	})

	t.Run("updates frozen habit by 3 days", func(t *testing.T) {
		habit := newHabit("Test", 2, 60)
		habit.CheckStep()
		habit.CheckStep()
		habit.Summary.CurrentStreak = 2
		habit.Summary.LongestStreak = 3
		habit.Freeze()
		habit.UpdateToPresent(3)

		isUpdated := isUpdated(&habit, 2, 3, TotalTime{}, [HistoryLen]Entry{nil, nil, nil, FrozenEntry{}, FrozenEntry{}, FrozenEntry{}})

		if !isUpdated {
			t.Error("habit has not been updated successfully")
		}
	})

	t.Run("updates habit by 10 days", func(t *testing.T) {
		habit := newHabit("Test", 2, 60)
		habit.CheckStep()
		habit.CheckStep()
		habit.UpdateToPresent(10)

		isUpdated := isUpdated(&habit, 0, 1, TotalTime{Hours: 2}, [HistoryLen]Entry{ActiveEntry{0, 2}, ActiveEntry{0, 2}, ActiveEntry{0, 2}, ActiveEntry{0, 2}, ActiveEntry{0, 2}, ActiveEntry{0, 2}})

		if !isUpdated {
			t.Error("habit has not been updated successfully")
		}
	})

	t.Run("updates frozen habit by 10 days", func(t *testing.T) {
		habit := newHabit("Test", 2, 60)
		habit.CheckStep()
		habit.CheckStep()
		habit.Summary.CurrentStreak = 2
		habit.Summary.LongestStreak = 3
		habit.Freeze()
		habit.UpdateToPresent(10)

		isUpdated := isUpdated(&habit, 2, 3, TotalTime{}, [HistoryLen]Entry{FrozenEntry{}, FrozenEntry{}, FrozenEntry{}, FrozenEntry{}, FrozenEntry{}, FrozenEntry{}})

		if !isUpdated {
			t.Error("habit has not been updated successfully")
		}
	})
}
