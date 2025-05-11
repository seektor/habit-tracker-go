package utils

import (
	"testing"
	"time"
)

func TestGetDaysDiff(t *testing.T) {
	var tests = []struct {
		from time.Time
		to   time.Time
		want int32
	}{
		{time.Date(2020, 11, 20, 12, 34, 0, 0, time.UTC),
			time.Date(2020, 11, 20, 15, 34, 0, 0, time.UTC),
			0},
		{time.Date(2020, 11, 20, 23, 59, 59, 0, time.UTC),
			time.Date(2020, 11, 21, 0, 0, 0, 0, time.UTC),
			1},
		{time.Date(2020, 11, 21, 12, 34, 0, 0, time.UTC),
			time.Date(2020, 11, 20, 0, 0, 0, 0, time.UTC),
			-1},
		{time.Date(2020, 11, 20, 12, 34, 0, 0, time.UTC),
			time.Date(2021, 11, 20, 15, 34, 0, 0, time.UTC),
			365},
	}

	for _, tt := range tests {

		t.Run("calculates days difference", func(t *testing.T) {
			diff := GetDaysDiff(tt.from, tt.to)

			if diff != tt.want {
				t.Errorf("invalid days difference, expected: %d, got: %d", tt.want, diff)
			}
		})
	}
}
