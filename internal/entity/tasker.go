package entity

import "errors"

var TaskTimeRangeNotFoundErr = errors.New("TimeRange not found in this task")
var TaskDurationImpossibleErr = errors.New("The current time is below the time of the time range")
var TaskTwoUnfinishedTimeRangeErr = errors.New("There are two unfinished time ranges")

func CreateTask(id DbId, name string, date Dater) Tasker {
	return &Task{
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
	GetName() ProperNoun
	GetId() DbId
}

type Task struct {
	id         DbId
	name       ProperNoun
	date       Dater
	timeRanges map[DbId]TimeRanger
}

var _ Tasker = &Task{}

func (t *Task) Rename(newName string) {
	t.name = CreateProperNoun(newName)
}

func (t Task) GetId() DbId {
	return t.id
}

func (t Task) GetName() ProperNoun {
	return t.name
}

func (t *Task) SetTimeRange(timeRange TimeRanger) {
	t.timeRanges[timeRange.GetId()] = timeRange
}

func (t Task) GetTimeRange(id DbId) (TimeRanger, error) {
	timeRange, found := t.timeRanges[id]

	if !found {
		return nil, TaskTimeRangeNotFoundErr
	}

	return timeRange, nil
}

func (t *Task) Duration() (Duration, error) {
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

func (t Task) IsRunning() bool {

	countNotRunning := 0

	for _, timeRange := range t.timeRanges {
		if !timeRange.HasEnded() {
			countNotRunning += 1
		}
	}

	return countNotRunning == 1
}
