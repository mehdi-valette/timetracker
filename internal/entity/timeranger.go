package entity

import (
	"errors"
)

var TimeRangeDurationNoEndErr = errors.New("no end value")
var TimeRangeDurationNoStartErr = errors.New("no start value")
var TimeRangeDurationEndBeforeStartErr = errors.New("end before start")

func CreateTimeRange(id DbId, date Dater) TimeRanger {
	return &timeRange{
		id:    id,
		date:  date,
		start: Timestamp{},
		end:   Timestamp{},
	}
}

type TimeRanger interface {
	Duration() (Duration, error)
	End()
	GetEnd() Timestamp
	GetId() DbId
	GetStart() Timestamp
	HasEnded() bool
	HasStated() bool
	Start()
}

type timeRange struct {
	id              DbId
	date            Dater
	rangeHasStarted bool
	rangeHasEnded   bool
	start           Timestamp
	end             Timestamp
}

var _ TimeRanger = &timeRange{}

func (tr timeRange) Duration() (Duration, error) {
	if !tr.rangeHasStarted {
		return 0, TimeRangeDurationNoStartErr
	}

	if !tr.rangeHasEnded {
		return 0, TimeRangeDurationNoEndErr
	}

	if tr.end.IsBefore(tr.start) {
		return 0, TimeRangeDurationEndBeforeStartErr
	}

	return Duration(tr.end.GetSeconds() - tr.start.GetSeconds()), nil
}

func (tr *timeRange) GetEnd() Timestamp {
	return tr.end
}

func (tr timeRange) GetId() DbId {
	return tr.id
}

func (tr *timeRange) GetStart() Timestamp {
	return tr.start
}

func (tr *timeRange) End() {
	tr.end = tr.date.Now()
	tr.rangeHasEnded = true
}

func (tr timeRange) HasEnded() bool {
	return tr.rangeHasEnded
}

func (tr timeRange) HasStated() bool {
	return tr.rangeHasStarted
}

func (tr *timeRange) Start() {
	tr.start = tr.date.Now()
	tr.rangeHasStarted = true
}
