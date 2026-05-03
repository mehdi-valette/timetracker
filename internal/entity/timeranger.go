package entity

import (
	"errors"
)

var TimeRangeDurationNoEndErr = errors.New("no end value")
var TimeRangeDurationNoStartErr = errors.New("no start value")
var TimeRangeDurationEndBeforeStartErr = errors.New("end before start")

func CreateTimeRange(id DbId, taskId DbId, date Dater) TimeRanger {
	return &timeRange{
		date:   date,
		end:    Timestamp{},
		id:     id,
		start:  Timestamp{},
		taskId: taskId,
	}
}

type TimeRanger interface {
	Duration() (Duration, error)
	End()
	GetEnd() Timestamp
	GetId() DbId
	GetStart() Timestamp
	GetTaskId() DbId
	HasEnded() bool
	HasStated() bool
	Start()
}

type timeRange struct {
	date            Dater
	end             Timestamp
	id              DbId
	rangeHasEnded   bool
	rangeHasStarted bool
	start           Timestamp
	taskId          DbId
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

func (tr timeRange) GetStart() Timestamp {
	return tr.start
}

func (tr timeRange) GetTaskId() DbId {
	return tr.taskId
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
