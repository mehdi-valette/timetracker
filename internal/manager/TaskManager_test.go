package manager

import (
	"errors"
	"testing"

	"github.com/mehdi-valette/timetracker/internal/entity"
)

type mockCalls struct {
	method string
	arg    any
}

type taskRepositoryMock struct {
	lastId    entity.DbId
	createErr error
	saverErr  error
	getErr    error
	calls     []mockCalls
	taskList  map[entity.DbId]entity.Tasker
}

var _ TaskPersister = &taskRepositoryMock{}

func createTaskRepositoryMock() *taskRepositoryMock {
	return &taskRepositoryMock{
		lastId:   0,
		calls:    make([]mockCalls, 0),
		taskList: make(map[entity.DbId]entity.Tasker),
	}
}

func (p *taskRepositoryMock) HasBeenCalledWith(method string, arg any) bool {
	return p.HasBeenCalledTimes(method, arg) > 0
}

func (p *taskRepositoryMock) HasBeenCalledTimes(method string, arg any) uint8 {
	times := uint8(0)

	for _, call := range p.calls {
		if call.method == method && arg == arg {
			times++
		}
	}

	return times
}

func (p *taskRepositoryMock) Create() (entity.DbId, error) {
	p.lastId += 1

	p.calls = append(p.calls, mockCalls{"Create", nil})

	p.taskList[p.lastId] = entity.CreateTask(p.lastId, "", DateMock{})

	return p.lastId, p.createErr
}

func (p *taskRepositoryMock) Save(task entity.Tasker) error {
	p.calls = append(p.calls, mockCalls{"Save", task})

	p.taskList[task.GetId()] = task

	return p.saverErr
}

func (p *taskRepositoryMock) Get(taskId entity.DbId) (entity.Tasker, error) {
	if p.getErr != nil {
		return &entity.Task{}, p.getErr
	}

	task, taskFound := p.taskList[taskId]

	if !taskFound {
		return &entity.Task{}, errors.New("task not found")
	}

	return task, nil
}

func (p *taskRepositoryMock) Delete(entity.DbId) error { panic("todo") }

func createTestTaskManager(persistMock *taskRepositoryMock) (TaskManager, TimeRangePersister) {
	timeRangeRepository := createTimeRangeRepositoryMock()

	taskManager := CreateTaskManager(persistMock, CreateTimeRangeManager(timeRangeRepository, DateMock{}), entity.CreateDate())

	return taskManager, timeRangeRepository
}

func createTestTaskManagerWithMock(persistMock *taskRepositoryMock, timeRangePersister TimeRangePersister) TaskManager {
	return CreateTaskManager(persistMock, CreateTimeRangeManager(timeRangePersister, DateMock{}), entity.CreateDate())
}

func TestTaskManagerCreate(t *testing.T) {
	taskRepositoryMock := createTaskRepositoryMock()
	manager, _ := createTestTaskManager(taskRepositoryMock)

	task, createError := manager.Create("   hello     ")

	if createError != nil {
		t.Error("should not return an error")
	}

	if !taskRepositoryMock.HasBeenCalledWith("Save", task) {
		t.Error("should have called \"save\"")
	}

	if taskRepositoryMock.HasBeenCalledTimes("Save", task) != 1 {
		t.Error("should have called \"save\" 1 time")
	}

	if task.GetName() != "hello" {
		t.Errorf("the name should be \"hello\n, it is %s", task.GetName())
	}

	if !taskRepositoryMock.HasBeenCalledWith("Create", nil) {
		t.Error("should have called create")
	}
}

func TestTaskManagerCreateError(t *testing.T) {
	persistMock := createTaskRepositoryMock()
	expectedErr := errors.New("create error")

	persistMock.createErr = expectedErr

	manager, _ := createTestTaskManager(persistMock)

	task, createError := manager.Create("   hello     ")

	if !errors.Is(createError, expectedErr) {
		t.Error("should return an error")
	}

	if task != nil {
		t.Error("task should be nil")
	}

	if !persistMock.HasBeenCalledWith("Create", nil) {
		t.Error("should have called create")
	}
}

func TestTaskManagerStart(t *testing.T) {
	persistMock := createTaskRepositoryMock()

	manager, _ := createTestTaskManager(persistMock)

	task, _ := manager.Create("hello")

	if task.IsRunning() {
		t.Error("the task should not yet be running")
	}

	persistMock.calls = []mockCalls{}

	manager.Start(task.GetId())

	if !task.IsRunning() {
		t.Error("the task should be running")
	}
}

func TestTaskManagerStartNotFound(t *testing.T) {
	expectedErr := errors.New("task not found")

	persistMock := createTaskRepositoryMock()
	persistMock.getErr = expectedErr

	manager, _ := createTestTaskManager(persistMock)

	task, _ := manager.Create("hello")

	if task.IsRunning() {
		t.Error("the task should not yet be running")
	}

	startErr := manager.Start(-1)

	if !errors.Is(startErr, expectedErr) {
		t.Error("should return an error")
	}
}

func TestTaskManagerStartAlreadyRunning(t *testing.T) {
	persistMock := createTaskRepositoryMock()

	manager, timeRangeRepository := createTestTaskManager(persistMock)

	task, _ := manager.Create("hello")

	if task.IsRunning() {
		t.Error("the task should not yet be running")
	}

	if len(timeRangeRepository.ListByTaskId(task.GetId())) != 0 {
		t.Error("task should have 0 time ranges")
	}

	firstStartErr := manager.Start(task.GetId())

	if firstStartErr != nil {
		t.Error("should not return an error")
	}

	if len(timeRangeRepository.ListByTaskId(task.GetId())) != 1 {
		t.Error("task should have 1 time range")
	}

	secondStartErr := manager.Start(task.GetId())

	if !errors.Is(secondStartErr, TaskManagerTaskRunningErr) {
		t.Error("should return an error")
	}

	if len(timeRangeRepository.ListByTaskId(task.GetId())) != 1 {
		t.Error("task should have 1 time range")
	}
}

func TestTaskManagerStartTimeRangeError(t *testing.T) {
	expectedError := errors.New("cannot create")

	persistMock := createTaskRepositoryMock()
	timeRangeRepoMock := createTimeRangeRepositoryMock()
	timeRangeRepoMock.createError = expectedError

	manager := createTestTaskManagerWithMock(persistMock, timeRangeRepoMock)

	task, _ := manager.Create("hello")

	startErr := manager.Start(task.GetId())

	if !errors.Is(startErr, expectedError) {
		t.Error("should return an error")
	}
}

func TestTaskManagerSave(t *testing.T) {
	persistMock := createTaskRepositoryMock()

	manager, _ := createTestTaskManager(persistMock)

	task, _ := manager.Create("my task")

	if !persistMock.HasBeenCalledWith("Save", task) {
		t.Error("should have called \"save\"")
	}

	if persistMock.HasBeenCalledTimes("Save", task) != 1 {
		t.Error("should have called \"save\" 1 time")
	}

	manager.Save(task.GetId())

	if persistMock.HasBeenCalledTimes("Save", task) != 2 {
		t.Error("should have called \"save\" 2 times")
	}
}

func TestTaskManagerSaveTaskNotFound(t *testing.T) {
	expectedErr := errors.New("task not found")
	taskRepositoryMock := createTaskRepositoryMock()

	manager, _ := createTestTaskManager(taskRepositoryMock)

	task, _ := manager.Create("my task")

	if !taskRepositoryMock.HasBeenCalledWith("Save", task) {
		t.Error("should have called \"save\"")
	}

	if taskRepositoryMock.HasBeenCalledTimes("Save", task) != 1 {
		t.Error("should have called \"save\" 1 time")
	}

	taskRepositoryMock.getErr = expectedErr

	saveErr := manager.Save(-1)

	if !errors.Is(saveErr, expectedErr) {
		t.Error("should have returned an error")
	}
}

func TestTaskManagerSaveRepositoryError(t *testing.T) {
	expectedError := errors.New("create error")

	taskRepositoryMock := createTaskRepositoryMock()

	manager, _ := createTestTaskManager(taskRepositoryMock)

	task, createErr := manager.Create("my task")

	if createErr != nil {
		t.Error("should not return an error")
	}

	taskRepositoryMock.saverErr = expectedError

	saveErr := manager.Save(task.GetId())

	if !errors.Is(saveErr, expectedError) {
		t.Error("should have returned an error")
	}
}
