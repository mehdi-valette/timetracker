package entity

import (
	"errors"
	"testing"
)

func TestTimeRangeDurationSuccess(t *testing.T) {
	end := CreateTimestampSeconds(1500)

	timerange := TimeRange{
		id:     1,
		taskId: 1,
		start:  CreateTimestampSeconds(1000),
		end:    &end,
	}

	expected := uint32(500)

	result, _ := timerange.Duration()

	if result != expected {
		t.Errorf("should be 500, got %d", expected)
	}
}

func TestTimeRangeDurationNoEnd(t *testing.T) {
	timerange := TimeRange{
		id:     1,
		taskId: 1,
		start:  CreateTimestampSeconds(1000),
		end:    nil,
	}

	_, err := timerange.Duration()

	if !errors.Is(err, TimeRangeDurationNoEndErr) {
		t.Errorf("shoudl return an error when there's no end")
	}
}

func TestTimeRangeDurationEndBeforeStart(t *testing.T) {
	end := CreateTimestampSeconds(500)

	timerange := TimeRange{
		id:     1,
		taskId: 1,
		start:  CreateTimestampSeconds(1000),
		end:    &end,
	}

	val, err := timerange.Duration()

	if val != 0 {
		t.Errorf("the return value should be 0")
	}

	if !errors.Is(err, TimeRangeDurationEndBeforeStartErr) {
		t.Errorf("shoudl return an error when there's no end")
	}
}

func TestTimeRangeDurationEndEqualsStart(t *testing.T) {
	end := CreateTimestampSeconds(1000)

	timerange := TimeRange{
		id:     1,
		taskId: 1,
		start:  CreateTimestampSeconds(1000),
		end:    &end,
	}

	val, err := timerange.Duration()

	if val != 0 {
		t.Errorf("when start = end the duration should be 0")
	}

	if err != nil {
		t.Errorf("should not return an error")
	}
}
