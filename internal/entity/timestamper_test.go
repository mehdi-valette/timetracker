package entity

import (
	"errors"
	"testing"
)

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

func TestTimestampTimeEllapsedSince(t *testing.T) {
	t1 := CreateTimestampSeconds(1000)
	t2 := CreateTimestampSeconds(500)

	ellapsed, err := t1.TimeEllapsedSince(t2)

	if err != nil {
		t.Errorf("should not return an error")
	}

	if ellapsed != 500 {
		t.Errorf("500 seconds should have ellapsed")
	}
}

func TestTimestampTimeEllapsedSinceZero(t *testing.T) {
	t1 := CreateTimestampSeconds(1000)
	t2 := CreateTimestampSeconds(1000)

	ellapsed, err := t1.TimeEllapsedSince(t2)

	if err != nil {
		t.Errorf("should not return an error")
	}

	if ellapsed != 0 {
		t.Errorf("0 seconds should have ellapsed")
	}
}

func TestTimestampTimeEllapsedSinceError(t *testing.T) {
	t1 := CreateTimestampSeconds(1000)
	t2 := CreateTimestampSeconds(2000)

	_, err := t1.TimeEllapsedSince(t2)

	if !errors.Is(err, TimestampErrorNegativeDurationErr) {
		t.Errorf("should return an error")
	}
}
