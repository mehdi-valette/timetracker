package entity

import (
	"errors"
	"testing"
	"time"
)

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

	timeRange := CreateTimeRange(12, 3, &dateMock)

	if timeRange.GetId() != 12 {
		t.Errorf("should have the ID set during its creation")
	}

	if timeRange.GetTaskId() != 3 {
		t.Errorf("should belong to a task from creation")
	}

	if timeRange.HasStated() || timeRange.HasEnded() {
		t.Errorf("should not have started or ended")
	}
}

func TestTimeRangeDurationSuccess(t *testing.T) {
	timerange := timeRange{
		id:              1,
		rangeHasStarted: true,
		rangeHasEnded:   true,
		start:           CreateTimestampSeconds(1000),
		end:             CreateTimestampSeconds(1500),
	}

	expected := Duration(500)

	result, _ := timerange.Duration()

	if result != expected {
		t.Errorf("should be 500, got %d", expected)
	}
}

func TestTimeRangeDurationNoStart(t *testing.T) {
	timerange := timeRange{
		id:    1,
		start: CreateTimestampSeconds(1000),
		end:   CreateTimestampSeconds(0),
	}

	_, err := timerange.Duration()

	if !errors.Is(err, TimeRangeDurationNoStartErr) {
		t.Errorf("shoudl return an error when there's no start")
	}
}

func TestTimeRangeDurationNoEnd(t *testing.T) {
	timeRange := CreateTimeRange(0, 0, date{})

	timeRange.Start()

	_, err := timeRange.Duration()

	if !errors.Is(err, TimeRangeDurationNoEndErr) {
		t.Errorf("shoudl return an error when there's no end")
	}
}

func TestTimeRangeDurationEndBeforeStart(t *testing.T) {
	timerange := timeRange{
		id:              1,
		rangeHasStarted: true,
		rangeHasEnded:   true,
		start:           CreateTimestampSeconds(1000),
		end:             CreateTimestampSeconds(500),
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
	timerange := timeRange{
		id:              1,
		rangeHasStarted: true,
		rangeHasEnded:   true,
		start:           CreateTimestampSeconds(1000),
		end:             CreateTimestampSeconds(1000),
	}

	val, err := timerange.Duration()

	if val != 0 {
		t.Errorf("when start = end the duration should be 0")
	}

	if err != nil {
		t.Errorf("should not return an error")
	}
}

func TestTimeRangeEnd(t *testing.T) {
	mockedDate := DaterMock{
		innerDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	timeRange := CreateTimeRange(0, 0, &mockedDate)

	if timeRange.HasStated() {
		t.Errorf("should not have started yet")
	}

	timeRange.Start()

	if !timeRange.HasStated() {
		t.Errorf("should have started")
	}

	if timeRange.HasEnded() {
		t.Errorf("should not have finished yet")
	}

	mockedDate.Set(time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC))

	timeRange.End()

	if !timeRange.HasEnded() {
		t.Errorf("should be finished")
	}

	duration, err := timeRange.Duration()

	if err != nil || duration != 3600 {
		t.Errorf("the time range should have a duration of 3600 seconds, but has %d, %e", duration, err)
	}
}

func TestTimeRangeGetId(t *testing.T) {
	chosenId := DbId(12)

	timeRange := CreateTimeRange(chosenId, 0, date{})

	if timeRange.GetId() != chosenId {
		t.Error("should return the correct ID")
	}
}

func TestTimeRangeGetStart(t *testing.T) {
	chosenDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	mockedDate := DaterMock{
		innerDate: chosenDate,
	}

	timeRange := CreateTimeRange(0, 0, mockedDate)

	timeRange.Start()

	if timeRange.GetStart() != Timestamp(chosenDate) {
		t.Error("should return the correct date")
	}
}

func TestTimeRangeGetEnd(t *testing.T) {
	startDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC)

	mockedDate := DaterMock{
		innerDate: startDate,
	}

	timeRange := CreateTimeRange(0, 0, &mockedDate)
	timeRange.Start()

	mockedDate.innerDate = endDate

	timeRange.End()

	if timeRange.GetEnd() != Timestamp(endDate) {
		t.Error("should return the correct date")
	}
}
