package entity

import "errors"

var TaskTimeRangeNotFoundErr = errors.New("TimeRange not found in this task")
var TaskDurationImpossibleErr = errors.New("The current time is below the time of the time range")
var TaskTwoUnfinishedTimeRangeErr = errors.New("There are two unfinished time ranges")

func CreateTask(id DbId, name string, date Dater) (Tasker, error) {
	noun, nounErr := CreateProperNoun(name)

	if nounErr != nil {
		return &Task{}, nounErr
	}

	return &Task{
		id:         id,
		name:       noun,
		date:       date,
		timeRanges: make([]TimeRanger, 0),
	}, nil
}

type Tasker interface {
	Rename(newName string) error
	SetTimeRange(timeRange TimeRanger)
	GetTimeRangeById(id DbId) (TimeRanger, error)
	GetLastTimeRange() (TimeRanger, error)
	Duration() (Duration, error)
	IsRunning() bool
	GetName() ProperNoun
	GetId() DbId
}

type Task struct {
	id         DbId
	name       ProperNoun
	date       Dater
	timeRanges []TimeRanger
}

var _ Tasker = &Task{}

func (t *Task) Rename(newName string) error {
	noun, nounErr := CreateProperNoun(newName)

	if nounErr != nil {
		return nounErr
	}

	t.name = noun

	return nil
}

func (t Task) GetId() DbId {
	return t.id
}

func (t Task) GetName() ProperNoun {
	return t.name
}

func (t *Task) SetTimeRange(timeRange TimeRanger) {
	t.timeRanges = append(t.timeRanges, timeRange)
}

func (t Task) GetTimeRangeById(timeRangeId DbId) (TimeRanger, error) {
	for _, timeRange := range t.timeRanges {
		if timeRange.GetId() == timeRangeId {
			return timeRange, nil
		}
	}

	return nil, TaskTimeRangeNotFoundErr
}

func (t Task) GetLastTimeRange() (TimeRanger, error) {
	lastIndex := len(t.timeRanges) - 1

	if lastIndex < 0 {
		return nil, TaskTimeRangeNotFoundErr
	}

	return t.timeRanges[lastIndex], nil
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
