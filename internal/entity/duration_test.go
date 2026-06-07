package entity

import "testing"

func TestDurationToSTring(t *testing.T) {
	type sample struct {
		duration Duration
		expected string
	}

	samples := []sample{
		{
			duration: Duration(4250),
			expected: "01:10:50",
		},
		{
			duration: Duration(10),
			expected: "00:00:10",
		},
		{
			duration: Duration(90),
			expected: "00:01:30",
		},
		{
			duration: Duration(0),
			expected: "00:00:00",
		},
	}

	for _, test := range samples {
		if test.duration.ToString() != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, test.duration.ToString())
		}
	}
}
