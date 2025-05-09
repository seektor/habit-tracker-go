package habit

import "testing"

func TestStringify(t *testing.T) {

	t.Run("Stringifies minutes", func(t *testing.T) {
		h := TotalTime{Minutes: 10}
		got := h.Stringify()
		want := "10 Minutes"

		if want != got {
			t.Errorf("Invalid time string, expected: %s, got: %s", want, got)
		}
	})

	t.Run("Stringifies hours and minutes", func(t *testing.T) {
		h := TotalTime{Minutes: 10, Hours: 20}
		got := h.Stringify()
		want := "20 Hours 10 Minutes"

		if want != got {
			t.Errorf("Invalid time string, expected: %s, got: %s", want, got)
		}
	})

	t.Run("Stringifies days, hours and minutes", func(t *testing.T) {
		h := TotalTime{Minutes: 10, Hours: 20, Days: 30}
		got := h.Stringify()
		want := "30 Days 20 Hours 10 Minutes"

		if want != got {
			t.Errorf("Invalid time string, expected: %s, got: %s", want, got)
		}
	})
}

func TestAdd(t *testing.T) {
	var tests = []struct {
		initial TotalTime
		minutes int16
		want    TotalTime
	}{
		{TotalTime{Minutes: 0}, 10, TotalTime{Minutes: 10}},
		{TotalTime{Minutes: 0}, 70, TotalTime{Minutes: 10, Hours: 1}},
		{TotalTime{Hours: 23, Minutes: 0}, 130, TotalTime{Days: 1, Hours: 1, Minutes: 10}},
	}

	for _, tt := range tests {

		t.Run("Adds minutes", func(t *testing.T) {
			tt.initial.Add(tt.minutes)

			if tt.initial != tt.want {
				t.Errorf("Invalid total time, expected: %s, got: %s", tt.want.Stringify(), tt.initial.Stringify())
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	var tests = []struct {
		initial TotalTime
		minutes int16
		want    TotalTime
	}{
		{TotalTime{Minutes: 0}, 10, TotalTime{Minutes: 0}},
		{TotalTime{Minutes: 20}, 10, TotalTime{Minutes: 10}},
		{TotalTime{Hours: 2}, 60, TotalTime{Hours: 1}},
		{TotalTime{Hours: 1, Minutes: 0}, 10, TotalTime{Hours: 0, Minutes: 50}},
		{TotalTime{Hours: 23, Minutes: 30}, 130, TotalTime{Hours: 21, Minutes: 20}},
		{TotalTime{Days: 2, Hours: 1, Minutes: 10}, 120, TotalTime{Days: 1, Hours: 23, Minutes: 10}},
		{TotalTime{Days: 2}, 60 * 24, TotalTime{Days: 1}},
	}

	for _, tt := range tests {

		t.Run("Subtracts minutes", func(t *testing.T) {
			tt.initial.Subtract(tt.minutes)

			if tt.initial != tt.want {
				t.Errorf("Invalid total time, expected: %s, got: %s", tt.want.Stringify(), tt.initial.Stringify())
			}
		})
	}
}
