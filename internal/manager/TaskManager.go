package manager

import "github.com/mehdi-valette/timetracker/internal/entity"

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

func CreateTaskManager(persister TaskPersister, date entity.Dater) TaskManager {
	return &TaskManagement{
		persister: persister,
		date:      date,
	}
}

type TaskManagement struct {
	persister TaskPersister
	date      entity.Dater
}

var _ TaskManager = &TaskManagement{}

func (tm *TaskManagement) Create(name string) (entity.Tasker, error) {
	taskId, createErr := tm.persister.Create()

	if createErr != nil {
		return nil, createErr
	}

	task := entity.CreateTask(taskId, name, tm.date)

	return task, nil
}

func (tm *TaskManagement) Start(id entity.DbId) error {
	panic("todo")
}

func (tm *TaskManagement) Save(id entity.DbId) error {
	panic("todo")
}

func (tm *TaskManagement) Stop(id entity.DbId) error {
	panic("todo")
}
