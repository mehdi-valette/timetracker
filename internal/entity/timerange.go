package entity

import "errors"

var TimeRangeDurationNoEndErr = errors.New("no end value")
var TimeRangeDurationEndBeforeStartErr = errors.New("end before start")

type TimeRange struct {
	id     DbId
	taskId DbId
	start  Timestamp
	end    *Timestamp
}

func (tr TimeRange) Duration() (uint32, error) {
	if tr.end == nil {
		return 0, TimeRangeDurationNoEndErr
	}

	if tr.end.IsBefore(tr.start) {
		return 0, TimeRangeDurationEndBeforeStartErr
	}

	return uint32(tr.end.GetSeconds() - tr.start.GetSeconds()), nil
}
