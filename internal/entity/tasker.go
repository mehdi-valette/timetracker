package entity

import "errors"

var TaskTimeRangeNotFoundErr = errors.New("TimeRange not found in this task")
var TaskDurationImpossibleErr = errors.New("The current time is below the time of the time range")
var TaskTwoUnfinishedTimeRangeErr = errors.New("There are two unfinished time ranges")

func CreateTask(id DbId, name string, date Dater) Tasker {
	return &task{
		id:         id,
		name:       CreateProperNoun(name),
		date:       date,
		timeRanges: make(map[DbId]TimeRanger),
	}
}

type Tasker interface {
	Rename(newName string)
	SetTimeRange(timeRange TimeRanger)
	GetTimeRange(id DbId) (TimeRanger, error)
	Duration() (Duration, error)
	IsRunning() bool
	Name() ProperNoun
	Id() DbId
}

type task struct {
	id         DbId
	name       ProperNoun
	date       Dater
	timeRanges map[DbId]TimeRanger
}

var _ Tasker = &task{}

func (t *task) Rename(newName string) {
	t.name = CreateProperNoun(newName)
}

func (t task) Id() DbId {
	return t.id
}

func (t task) Name() ProperNoun {
	return t.name
}

func (t *task) SetTimeRange(timeRange TimeRanger) {
	t.timeRanges[timeRange.GetId()] = timeRange
}

func (t task) GetTimeRange(id DbId) (TimeRanger, error) {
	timeRange, found := t.timeRanges[id]

	if !found {
		return nil, TaskTimeRangeNotFoundErr
	}

	return timeRange, nil
}

func (t *task) Duration() (Duration, error) {
	totalDuration := Duration(0)

	firstUnfinishedtimeRange := true

	for _, timeRange := range t.timeRanges {
		duration, err := timeRange.Duration()

		if err != nil {
			if firstUnfinishedtimeRange {
				firstUnfinishedtimeRange = false
			} else {
				return 0, TaskTwoUnfinishedTimeRangeErr
			}

			ellapsed, err := t.date.Now().TimeEllapsedSince(timeRange.GetStart())

			if err != nil {
				return 0, TaskDurationImpossibleErr
			}

			totalDuration += ellapsed

			continue
		}

		totalDuration += duration
	}

	return totalDuration, nil
}

func (t task) IsRunning() bool {

	countNotRunning := 0

	for _, timeRange := range t.timeRanges {
		if !timeRange.HasEnded() {
			countNotRunning += 1
		}
	}

	return countNotRunning == 1
}
