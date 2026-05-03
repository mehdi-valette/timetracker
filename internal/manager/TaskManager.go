package manager

import (
	"errors"

	"github.com/mehdi-valette/timetracker/internal/entity"
)

var TaskNotFoundErr = errors.New("task not found")

type TaskPersister interface {
	Create() (entity.DbId, error)
	Save(entity.Tasker) error
	Delete(entity.DbId) error
}

type TaskManager interface {
	Create(name string) (entity.Tasker, error)
	Start(id entity.DbId) error
	Save(id entity.DbId) error
	Stop(id entity.DbId) error
}

func CreateTaskManager(persister TaskPersister, timeRangeManager TimeRangeManager, date entity.Dater) TaskManager {
	return &TaskManagement{
		taskRepository:   persister,
		date:             date,
		timeRangeManager: timeRangeManager,

		tasks: make(map[entity.DbId]entity.Tasker),
	}
}

type TaskManagement struct {
	taskRepository   TaskPersister
	timeRangeManager TimeRangeManager
	date             entity.Dater

	tasks map[entity.DbId]entity.Tasker
}

var _ TaskManager = &TaskManagement{}

func (tm *TaskManagement) Create(name string) (entity.Tasker, error) {
	taskId, createErr := tm.taskRepository.Create()

	if createErr != nil {
		return nil, createErr
	}

	task := entity.CreateTask(taskId, name, tm.date)

	tm.tasks[task.GetId()] = task

	return task, nil
}

func (tm *TaskManagement) Start(id entity.DbId) error {
	task, taskFound := tm.tasks[id]

	if !taskFound {
		return TaskNotFoundErr
	}

	if task.IsRunning() {
		return nil
	}

	timeRange, timeRangeErr := tm.timeRangeManager.Create(task.GetId())

	if timeRangeErr != nil {
		return timeRangeErr
	}

	task.SetTimeRange(timeRange)

	tm.timeRangeManager.Start(timeRange.GetId())

	return nil
}

func (tm *TaskManagement) Save(id entity.DbId) error {
	panic("todo")
}

func (tm *TaskManagement) Stop(id entity.DbId) error {
	panic("todo")
}
