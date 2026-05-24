package manager

import (
	"errors"

	"github.com/mehdi-valette/timetracker/internal/entity"
)

var TimeRangeManagerTimeRangeNotFoundErr = errors.New("cannot find the time range")

type TimeRangePersister interface {
	Create(taskId entity.DbId) (entity.DbId, error)
	Save(timeRange entity.TimeRanger) error
	Get(timeRangeId entity.DbId) (entity.TimeRanger, error)
	Delete(id entity.DbId) error
	ListByTaskId(taskId entity.DbId) ([]entity.TimeRanger, error)
}

type TimeRangeManager interface {
	Create(taskId entity.DbId) (entity.TimeRanger, error)
	Start(id entity.DbId) error
	Stop(id entity.DbId) error
	Save(timeRange entity.TimeRanger) error
	ListByTaskId(taskId entity.DbId) ([]entity.TimeRanger, error)
}

func CreateTimeRangeManager(repo TimeRangePersister, date entity.Dater) TimeRangeManager {
	return &TimeRangeManagement{
		repository: repo,
		date:       date,
	}
}

type TimeRangeManagement struct {
	repository TimeRangePersister
	date       entity.Dater
}

func (trm *TimeRangeManagement) Create(taskId entity.DbId) (entity.TimeRanger, error) {
	id, createError := trm.repository.Create(taskId)

	if createError != nil {
		return nil, createError
	}

	timeRange := entity.CreateTimeRange(id, taskId, trm.date)

	trm.repository.Save(timeRange)

	return timeRange, createError
}

func (trm *TimeRangeManagement) Start(id entity.DbId) error {
	timeRange, getErr := trm.repository.Get(id)

	if getErr != nil {
		return getErr
	}

	timeRange.Start()

	return trm.repository.Save(timeRange)
}

func (trm *TimeRangeManagement) Stop(id entity.DbId) error {
	timeRange, getErr := trm.repository.Get(id)

	if getErr != nil {
		return TimeRangeManagerTimeRangeNotFoundErr
	}

	timeRange.End()

	return trm.repository.Save(timeRange)
}

func (trm *TimeRangeManagement) Delete(id entity.DbId) error {
	return trm.repository.Delete(id)
}

func (trm *TimeRangeManagement) Save(timeRange entity.TimeRanger) error {
	return trm.repository.Save(timeRange)
}

// ListByTaskId implements [TimeRangeManager].
func (trm *TimeRangeManagement) ListByTaskId(taskId entity.DbId) ([]entity.TimeRanger, error) {
	return trm.repository.ListByTaskId(taskId)
}

var _ TimeRangeManager = &TimeRangeManagement{}
