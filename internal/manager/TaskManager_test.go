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

type persisterMock struct {
	lastId    entity.DbId
	createErr error
	saverErr  error
	calls     []mockCalls
}

func createPersisterMock() *persisterMock {
	return &persisterMock{
		lastId: 0,
		calls:  make([]mockCalls, 0),
	}
}

func (p *persisterMock) HasBeenCalledWith(method string, arg any) bool {
	return p.HasBeenCalledTimes(method, arg) > 0
}

func (p *persisterMock) HasBeenCalledTimes(method string, arg any) uint8 {
	times := uint8(0)

	for _, call := range p.calls {
		if call.method == method && arg == arg {
			times++
		}
	}

	return times
}

func (p *persisterMock) Create() (entity.DbId, error) {
	p.lastId += 1

	p.calls = append(p.calls, mockCalls{"Create", nil})

	return p.lastId, p.createErr
}

func (p *persisterMock) Save(task entity.Tasker) error {
	p.calls = append(p.calls, mockCalls{"Save", task})

	return p.saverErr
}

func (p *persisterMock) Delete(entity.DbId) error { panic("todo") }

func createTestTaskManager(persistMock *persisterMock) (TaskManager, TimeRangePersister) {
	timeRangeRepository := createTimeRangeRepositoryMock()

	taskManager := CreateTaskManager(persistMock, CreateTimeRangeManager(timeRangeRepository, DateMock{}), entity.CreateDate())

	return taskManager, timeRangeRepository
}

func createTestTaskManagerWithMock(persistMock *persisterMock, timeRangePersister TimeRangePersister) TaskManager {
	return CreateTaskManager(persistMock, CreateTimeRangeManager(timeRangePersister, DateMock{}), entity.CreateDate())
}

func TestTaskManagerCreate(t *testing.T) {
	persistMock := createPersisterMock()
	manager, _ := createTestTaskManager(persistMock)

	task, createError := manager.Create("   hello     ")

	if createError != nil {
		t.Error("should not return an error")
	}

	if !persistMock.HasBeenCalledWith("Save", task) {
		t.Error("should have called \"save\"")
	}

	if persistMock.HasBeenCalledTimes("Save", task) != 1 {
		t.Error("should have called \"save\" 1 time")
	}

	if task.GetName() != "hello" {
		t.Errorf("the name should be \"hello\n, it is %s", task.GetName())
	}

	if !persistMock.HasBeenCalledWith("Create", nil) {
		t.Error("should have called create")
	}
}

func TestTaskManagerCreateError(t *testing.T) {
	persistMock := createPersisterMock()
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
	persistMock := createPersisterMock()

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
	persistMock := createPersisterMock()

	manager, _ := createTestTaskManager(persistMock)

	task, _ := manager.Create("hello")

	if task.IsRunning() {
		t.Error("the task should not yet be running")
	}

	startErr := manager.Start(-1)

	if !errors.Is(startErr, TaskManagerTaskNotFoundErr) {
		t.Error("should return an error")
	}
}

func TestTaskManagerStartAlreadyRunning(t *testing.T) {
	persistMock := createPersisterMock()

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

	persistMock := createPersisterMock()
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
	persisterMock := createPersisterMock()

	manager, _ := createTestTaskManager(persisterMock)

	task, _ := manager.Create("my task")

	if !persisterMock.HasBeenCalledWith("Save", task) {
		t.Error("should have called \"save\"")
	}

	if persisterMock.HasBeenCalledTimes("Save", task) != 1 {
		t.Error("should have called \"save\" 1 time")
	}

	manager.Save(task.GetId())

	if persisterMock.HasBeenCalledTimes("Save", task) != 2 {
		t.Error("should have called \"save\" 2 times")
	}
}

func TestTaskManagerSaveTaskNotFound(t *testing.T) {
	persisterMock := createPersisterMock()

	manager, _ := createTestTaskManager(persisterMock)

	task, _ := manager.Create("my task")

	if !persisterMock.HasBeenCalledWith("Save", task) {
		t.Error("should have called \"save\"")
	}

	if persisterMock.HasBeenCalledTimes("Save", task) != 1 {
		t.Error("should have called \"save\" 1 time")
	}

	saveErr := manager.Save(-1)

	if !errors.Is(saveErr, TaskManagerTaskNotFoundErr) {
		t.Error("should have returned an error")
	}
}

func TestTaskManagerSaveRepositoryError(t *testing.T) {
	expectedError := errors.New("create error")

	persisterMock := createPersisterMock()

	manager, _ := createTestTaskManager(persisterMock)

	task, createErr := manager.Create("my task")

	if createErr != nil {
		t.Error("should not return an error")
	}

	persisterMock.saverErr = expectedError

	saveErr := manager.Save(task.GetId())

	if !errors.Is(saveErr, expectedError) {
		t.Error("should have returned an error")
	}
}
