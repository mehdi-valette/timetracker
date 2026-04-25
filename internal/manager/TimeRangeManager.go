package manager

import (
	"errors"

	"github.com/mehdi-valette/timetracker/internal/entity"
)

var TimeRangeManagerTimeRangeNotFoundErr = errors.New("cannot find the time range")

type TimeRangeRepositoryManager interface {
	Create() (entity.DbId, error)
	Save(timeRange entity.TimeRanger) error
	Delete(id entity.DbId) error
}

type TimeRangeManager interface {
	Create() (entity.TimeRanger, error)
	Start(id entity.DbId) error
	Stop(id entity.DbId) error
	Save(id entity.DbId) error
}

type TimeRangeManagement struct {
	repository TimeRangeRepositoryManager
	timeRanges map[entity.DbId]entity.TimeRanger
	date       entity.Dater

	// timeRange ID => task ID
	timeRangeToTask map[entity.DbId]entity.DbId
}

var _ TimeRangeManager = &TimeRangeManagement{}

func CreateTimeRangeManager(repo TimeRangeRepositoryManager, date entity.Dater) TimeRangeManager {
	return &TimeRangeManagement{
		repository:      repo,
		date:            date,
		timeRanges:      make(map[entity.DbId]entity.TimeRanger),
		timeRangeToTask: make(map[entity.DbId]entity.DbId),
	}
}

func (trm *TimeRangeManagement) Create() (entity.TimeRanger, error) {
	id, createError := trm.repository.Create()

	if createError != nil {
		return nil, createError
	}

	timeRange := entity.CreateTimeRange(id, trm.date)

	trm.timeRanges[timeRange.GetId()] = timeRange

	trm.repository.Save(timeRange)

	return timeRange, createError
}

func (trm *TimeRangeManagement) Start(id entity.DbId) error {
	timeRange, found := trm.timeRanges[id]

	if !found {
		return TimeRangeManagerTimeRangeNotFoundErr
	}

	timeRange.Start()

	return trm.repository.Save(timeRange)
}

func (trm *TimeRangeManagement) Stop(id entity.DbId) error {
	timeRange, found := trm.timeRanges[id]

	if !found {
		return TimeRangeManagerTimeRangeNotFoundErr
	}

	timeRange.End()

	return trm.repository.Save(timeRange)
}

func (trm *TimeRangeManagement) Delete(id entity.DbId) error {
	_, found := trm.timeRanges[id]

	if !found {
		return TimeRangeManagerTimeRangeNotFoundErr
	}

	deleteErr := trm.repository.Delete(id)

	if deleteErr != nil {
		return deleteErr
	}

	delete(trm.timeRanges, id)

	return nil
}

func (trm *TimeRangeManagement) Save(id entity.DbId) error {
	timeRange, found := trm.timeRanges[id]

	if !found {
		return TimeRangeManagerTimeRangeNotFoundErr
	}

	return trm.repository.Save(timeRange)
}
