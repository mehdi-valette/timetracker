package entity

import (
	"errors"
)

var TimeRangeDurationNoEndErr = errors.New("no end value")
var TimeRangeDurationEndBeforeStartErr = errors.New("end before start")

type TimeRange struct {
	id          DbId
	date        Dater
	rangeHasEnd bool
	start       Timestamp
	end         Timestamp
}

func CreateTimeRange(id DbId, date Dater) TimeRange {
	return TimeRange{
		id:    id,
		date:  date,
		start: Timestamp(date.Now()),
		end:   Timestamp{},
	}
}

func (tr TimeRange) Duration() (Duration, error) {
	if !tr.rangeHasEnd {
		return 0, TimeRangeDurationNoEndErr
	}

	if tr.end.IsBefore(tr.start) {
		return 0, TimeRangeDurationEndBeforeStartErr
	}

	return Duration(tr.end.GetSeconds() - tr.start.GetSeconds()), nil
}

func (tr TimeRange) IsFinished() bool {
	return tr.rangeHasEnd
}

func (tr *TimeRange) End() {
	tr.end = tr.date.Now()
	tr.rangeHasEnd = true
}
