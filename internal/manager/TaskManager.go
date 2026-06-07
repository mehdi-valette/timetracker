package manager

import (
	"errors"

	"github.com/mehdi-valette/timetracker/internal/entity"
)

var TaskManagerTaskRunningErr = errors.New("task already running")
var TaskManagerTaskStoppedErr = errors.New("task already running")

type TaskPersister interface {
	Create() (entity.DbId, error)
	Save(entity.Tasker) error
	Delete(entity.DbId) error
	Get(entity.DbId) (entity.Tasker, error)
	List() ([]entity.Tasker, error)
}

type TaskManager interface {
	Create(name string) (entity.Tasker, error)
	Start(id entity.DbId) (entity.Tasker, error)
	Save(id entity.DbId) error
	Stop(id entity.DbId) error
	Get(id entity.DbId) (entity.Tasker, error)
	List() ([]entity.Tasker, error)
}

func CreateTaskManager(persister TaskPersister, timeRangeManager TimeRangeManager, date entity.Dater) TaskManager {
	return &TaskManagement{
		taskRepository:   persister,
		date:             date,
		timeRangeManager: timeRangeManager,
	}
}

type TaskManagement struct {
	taskRepository   TaskPersister
	timeRangeManager TimeRangeManager
	date             entity.Dater
}

// Get implements [TaskManager].
func (tm *TaskManagement) Get(taskId entity.DbId) (entity.Tasker, error) {
	timeRanges, trErr := tm.timeRangeManager.ListByTaskId(taskId)

	if trErr != nil {
		return nil, trErr
	}

	task, getErr := tm.taskRepository.Get(taskId)

	if getErr != nil {
		return nil, getErr
	}

	for _, timeRange := range timeRanges {
		task.SetTimeRange(timeRange)
	}

	return task, nil
}

// List implements [TaskManager].
func (tm *TaskManagement) List() ([]entity.Tasker, error) {
	return tm.taskRepository.List()
}

func (tm *TaskManagement) Create(name string) (entity.Tasker, error) {
	taskId, createErr := tm.taskRepository.Create()

	if createErr != nil {
		return nil, createErr
	}

	task := entity.CreateTask(taskId, name, tm.date)

	if saveErr := tm.taskRepository.Save(task); saveErr != nil {
		deleteErr := tm.taskRepository.Delete(taskId)

		return nil, errors.Join(saveErr, deleteErr)
	}

	return task, nil
}

func (tm *TaskManagement) Start(taskId entity.DbId) (entity.Tasker, error) {
	task, getErr := tm.Get(taskId)

	if getErr != nil {
		return nil, getErr
	}

	if task.IsRunning() {
		return nil, TaskManagerTaskRunningErr
	}

	timeRange, timeRangeErr := tm.timeRangeManager.Create(task.GetId())

	if timeRangeErr != nil {
		return nil, timeRangeErr
	}

	timeRange, startErr := tm.timeRangeManager.Start(timeRange.GetId())

	if startErr != nil {
		return nil, startErr
	}

	task.SetTimeRange(timeRange)

	return task, nil
}

func (tm *TaskManagement) Save(taskId entity.DbId) error {
	task, getErr := tm.taskRepository.Get(taskId)

	if getErr != nil {
		return getErr
	}

	return tm.taskRepository.Save(task)
}

func (tm *TaskManagement) Stop(taskId entity.DbId) error {
	task, repoErr := tm.Get(taskId)

	if repoErr != nil {
		return repoErr
	}

	timeRange, getErr := task.GetLastTimeRange()

	if errors.Is(getErr, entity.TaskTimeRangeNotFoundErr) {
		return nil
	}

	stopErr := tm.timeRangeManager.Stop(timeRange.GetId())

	if stopErr != nil {
		return stopErr
	}

	return nil
}

var _ TaskManager = &TaskManagement{}
