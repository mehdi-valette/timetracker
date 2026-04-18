package entity

import (
	"errors"
	"testing"
	"time"
)

func TestTimeRangeDurationSuccess(t *testing.T) {
	timerange := TimeRange{
		id:          1,
		rangeHasEnd: true,
		start:       CreateTimestampSeconds(1000),
		end:         CreateTimestampSeconds(1500),
	}

	expected := Duration(500)

	result, _ := timerange.Duration()

	if result != expected {
		t.Errorf("should be 500, got %d", expected)
	}
}

func TestTimeRangeDurationNoEnd(t *testing.T) {
	timerange := TimeRange{
		id:    1,
		start: CreateTimestampSeconds(1000),
		end:   CreateTimestampSeconds(0),
	}

	_, err := timerange.Duration()

	if !errors.Is(err, TimeRangeDurationNoEndErr) {
		t.Errorf("shoudl return an error when there's no end")
	}
}

func TestTimeRangeDurationEndBeforeStart(t *testing.T) {
	timerange := TimeRange{
		id:          1,
		rangeHasEnd: true,
		start:       CreateTimestampSeconds(1000),
		end:         CreateTimestampSeconds(500),
	}

	val, err := timerange.Duration()

	if val != 0 {
		t.Errorf("the return value should be 0")
	}

	if !errors.Is(err, TimeRangeDurationEndBeforeStartErr) {
		t.Errorf("shoudl return an error when the end is before the start")
	}
}

func TestTimeRangeDurationEndEqualsStart(t *testing.T) {
	timerange := TimeRange{
		id:          1,
		rangeHasEnd: true,
		start:       CreateTimestampSeconds(1000),
		end:         CreateTimestampSeconds(1000),
	}

	val, err := timerange.Duration()

	if val != 0 {
		t.Errorf("when start = end the duration should be 0")
	}

	if err != nil {
		t.Errorf("should not return an error")
	}
}

type DaterMock struct {
	innerDate time.Time
}

var _ Dater = DaterMock{}

func (d DaterMock) Now() Timestamp {
	return Timestamp(d.innerDate)
}

func (d *DaterMock) Set(date time.Time) {
	d.innerDate = date
}

func TestCreateTimeRange(t *testing.T) {
	dateMock := DaterMock{
		innerDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	timeRange := CreateTimeRange(12, &dateMock)

	if timeRange.id != 12 {
		t.Errorf("should have the ID set during its creation")
	}

	if dateMock.Now() != timeRange.start {
		t.Errorf("should return the time set by the Dater parameter")
	}
}

func TestTimeRangeEnd(t *testing.T) {
	mockedDate := DaterMock{
		innerDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	timeRange := CreateTimeRange(0, &mockedDate)

	if timeRange.IsFinished() {
		t.Errorf("should not yet be finished")
	}

	mockedDate.Set(time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC))

	timeRange.End()

	if !timeRange.IsFinished() {
		t.Errorf("should be finished")
	}

	duration, err := timeRange.Duration()

	if err != nil || duration != 3600 {
		t.Errorf("the time range should have a duration of 3600 seconds, but has %d, %e", duration, err)
	}
}
