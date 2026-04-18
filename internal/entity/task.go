package entity

import "errors"

var TaskTimeRangeNotFoundErr = errors.New("TimeRange not found in this task")
var TaskDurationImpossibleErr = errors.New("The current time is below the time of the time range")
var TaskTwoUnfinishedTimeRangeErr = errors.New("There are two unfinished time ranges")

type Task struct {
	id         DbId
	name       Name
	date       Dater
	timeRanges map[DbId]*TimeRange
}

func CreateTask(id DbId, name string, date Dater) Task {
	return Task{
		id:         id,
		name:       CreateName(name),
		date:       date,
		timeRanges: make(map[DbId]*TimeRange),
	}
}

func (t *Task) Rename(newName string) {
	t.name = CreateName(newName)
}

func (t *Task) SetTimeRange(timeRange *TimeRange) {
	t.timeRanges[timeRange.id] = timeRange
}

func (t Task) GetTimeRange(id DbId) (*TimeRange, error) {
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

			ellapsed, err := t.date.Now().TimeEllapsedSince(timeRange.start)

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
		if !timeRange.IsFinished() {
			countNotRunning += 1
		}
	}

	return countNotRunning == 1
}
