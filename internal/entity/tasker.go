package entity

import (
	"errors"
	"fmt"
)

var TaskTimeRangeNotFoundErr = errors.New("TimeRange not found in this task")
var TaskDurationImpossibleErr = errors.New("The current time is below the time of the time range")
var TaskTwoUnfinishedTimeRangeErr = errors.New("There are two unfinished time ranges")

func CreateTask(id DbId, shortName string, name string, date Dater) Tasker {
	return &Task{
		id:         id,
		shortName:  CreateShortName(shortName),
		name:       CreateProperNoun(name),
		date:       date,
		timeRanges: make([]TimeRanger, 0),
	}
}

type Tasker interface {
	Rename(newName string)
	Reshortname(newShortName string)
	SetTimeRange(timeRange TimeRanger)
	GetTimeRangeById(id DbId) (TimeRanger, error)
	GetLastTimeRange() (TimeRanger, error)
	Duration() (Duration, error)
	IsRunning() bool
	GetName() ProperNoun
	GetShortName() ShortName
	GetId() DbId
	String() string
	StringShort() string
}

type Task struct {
	id         DbId
	shortName  ShortName
	name       ProperNoun
	date       Dater
	timeRanges []TimeRanger
}

var _ Tasker = &Task{}

// StringShort implements [Tasker].
func (t *Task) StringShort() string {
	return fmt.Sprintf("[%s] %s", t.GetShortName(), t.GetName())
}

// String implements [Tasker].
func (t *Task) String() string {
	duration, _ := t.Duration()

	indicator := " "
	if t.IsRunning() {
		indicator = "*"
	}

	return fmt.Sprintf("%s%s [%s] %s", indicator, duration.ToString(), t.GetShortName(), t.GetName())
}

// Reshortname implements [Tasker].
func (t *Task) Reshortname(newShortName string) {
	t.shortName = CreateShortName(newShortName)
}

// GetShortName implements [Tasker].
func (t *Task) GetShortName() ShortName {
	return t.shortName
}

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
