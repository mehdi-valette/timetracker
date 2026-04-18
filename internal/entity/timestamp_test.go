package entity

import "testing"

func TestTimestampCreate(t *testing.T) {
	expected := int64(1000)

	timestamp := CreateTimestampSeconds(expected)

	result := timestamp.GetSeconds()

	if result != expected {
		t.Errorf("expected %d seconds, got %d", expected, result)
	}
}

func TestTimestampIsBefore(t *testing.T) {
	type TestCase struct {
		ts1      Timestamp
		ts2      Timestamp
		expected bool
	}

	cases := []TestCase{
		{
			CreateTimestampSeconds(1000),
			CreateTimestampSeconds(2000),
			true,
		},
		{
			CreateTimestampSeconds(1000),
			CreateTimestampSeconds(1000),
			false,
		},
		{
			CreateTimestampSeconds(2000),
			CreateTimestampSeconds(1000),
			false,
		},
	}

	for _, c := range cases {
		if c.ts1.IsBefore(c.ts2) != c.expected {
			t.Errorf("is %d before %d should be %t", c.ts1.GetSeconds(), c.ts2.GetSeconds(), c.expected)
		}
	}
}
