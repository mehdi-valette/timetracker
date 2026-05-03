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
	calls     []mockCalls
}

func createPersisterMock() *persisterMock {
	return &persisterMock{
		lastId: 0,
		calls:  make([]mockCalls, 0),
	}
}

func (p *persisterMock) HasBeenCalledWith(method string, arg any) bool {
	for _, call := range p.calls {
		if call.method == method && arg == arg {
			return true
		}
	}

	return false
}

func (p *persisterMock) Create() (entity.DbId, error) {
	p.lastId += 1

	p.calls = append(p.calls, mockCalls{"Create", nil})

	return p.lastId, p.createErr
}

func (p *persisterMock) Save(entity.Tasker) error { panic("todo") }

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

	if !errors.Is(startErr, TaskNotFoundErr) {
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

	if secondStartErr != nil {
		t.Error("should not return an error")
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
