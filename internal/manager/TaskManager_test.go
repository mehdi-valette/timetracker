package manager

import (
	"errors"
	"testing"

	"github.com/mehdi-valette/timetracker/internal/entity"
	"github.com/mehdi-valette/timetracker/internal/test"
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
	listErr   error
	deleteErr error
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

	p.taskList[p.lastId], _ = entity.CreateTask(p.lastId, "", DateMock{})

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

func (p *taskRepositoryMock) Delete(taskId entity.DbId) error {
	p.calls = append(p.calls, mockCalls{"Delete", taskId})

	if p.deleteErr != nil {
		return p.deleteErr
	}

	delete(p.taskList, taskId)

	return nil
}

func (p *taskRepositoryMock) List() ([]entity.Tasker, error) {
	if p.listErr != nil {
		return []entity.Tasker{}, p.listErr
	}

	taskArray := make([]entity.Tasker, 0, len(p.taskList))

	for _, task := range p.taskList {
		taskArray = append(taskArray, task)
	}

	return taskArray, nil
}

func createTestTaskManager(persistMock *taskRepositoryMock) (TaskManager, *timeRangeRepositoryMock) {
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

func TestTaskManagerCreateEmptyName(t *testing.T) {
	taskRepositoryMock := createTaskRepositoryMock()
	manager, _ := createTestTaskManager(taskRepositoryMock)

	task, createError := manager.Create("        ")

	if !errors.Is(createError, entity.EmptyNounErr) {
		t.Error("should return an error")
	}

	if taskRepositoryMock.HasBeenCalledTimes("Save", task) != 0 {
		t.Error("should not have called \"save\"")
	}

	if taskRepositoryMock.HasBeenCalledTimes("Create", task) != 1 {
		t.Error("should have called create 1 time")
	}

	if taskRepositoryMock.HasBeenCalledTimes("Delete", task.GetId()) != 1 {
		t.Error("should have called delete 1 time")
	}
}

func TestTaskManagerCreateEmptyNameDeleteErr(t *testing.T) {
	expectedErr := errors.New("delete error")
	taskRepositoryMock := createTaskRepositoryMock()
	manager, _ := createTestTaskManager(taskRepositoryMock)
	taskRepositoryMock.deleteErr = expectedErr

	task, createError := manager.Create("        ")

	if !errors.Is(createError, entity.EmptyNounErr) {
		t.Error("should return the noun error")
	}

	if !errors.Is(createError, expectedErr) {
		t.Error("should return the delete error")
	}

	if taskRepositoryMock.HasBeenCalledTimes("Save", task) != 0 {
		t.Error("should not have called \"save\"")
	}

	if taskRepositoryMock.HasBeenCalledTimes("Create", task) != 1 {
		t.Error("should have called create 1 time")
	}

	if taskRepositoryMock.HasBeenCalledTimes("Delete", task.GetId()) != 1 {
		t.Error("should have called delete 1 time")
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

	timeRangesBeforeStart, _ := timeRangeRepository.ListByTaskId(task.GetId())

	firstStartErr := manager.Start(task.GetId())

	if firstStartErr != nil {
		t.Error("should not return an error")
	}

	timeRangesAfterFirstStart, _ := timeRangeRepository.ListByTaskId(task.GetId())

	secondStartErr := manager.Start(task.GetId())

	timeRangesAfterSecondStart, _ := timeRangeRepository.ListByTaskId(task.GetId())

	if len(timeRangesBeforeStart) != 0 {
		t.Error("task should have 0 time ranges")
	}

	if len(timeRangesAfterFirstStart) != 1 {
		t.Error("task should have 1 time range")
	}

	if !errors.Is(secondStartErr, TaskManagerTaskRunningErr) {
		t.Error("should return an error")
	}

	if len(timeRangesAfterSecondStart) != 1 {
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

func TestTaskManagerStop(t *testing.T) {
	taskRepositoryMock := createTaskRepositoryMock()

	manager, _ := createTestTaskManager(taskRepositoryMock)

	task, _ := manager.Create("my task")

	manager.Start(task.GetId())

	if !task.IsRunning() {
		t.Error("the task should be running")
	}

	stopErr := manager.Stop(task.GetId())

	if stopErr != nil {
		t.Error("should not return an error")
	}

	if task.IsRunning() {
		t.Error("the task should not be running")
	}
}

func TestTaskManagerStopNotRunning(t *testing.T) {
	taskRepositoryMock := createTaskRepositoryMock()

	taskManager, _ := createTestTaskManager(taskRepositoryMock)

	task, _ := taskManager.Create("my task")

	taskManager.Start(task.GetId())

	if !task.IsRunning() {
		t.Error("the task should be running")
	}

	taskManager.Stop(task.GetId())

	if task.IsRunning() {
		t.Error("the task should not be running")
	}

	stopErr := taskManager.Stop(task.GetId())

	if stopErr != nil {
		t.Error("should not return an error")
	}

	if task.IsRunning() {
		t.Error("the task should not be running")
	}
}

func TestTaskManagerStopNeverBegan(t *testing.T) {
	taskRepositoryMock := createTaskRepositoryMock()

	taskManager, _ := createTestTaskManager(taskRepositoryMock)

	task, _ := taskManager.Create("my task")

	if task.IsRunning() {
		t.Error("the task should not be running")
	}

	stopErr := taskManager.Stop(task.GetId())

	if stopErr != nil {
		t.Error("should not return an error")
	}

	if task.IsRunning() {
		t.Error("the task should not be running")
	}
}

func TestTaskManagerStopRepositoryError(t *testing.T) {
	expectedErr := errors.New("task not found")
	taskRepositoryMock := createTaskRepositoryMock()
	taskRepositoryMock.getErr = expectedErr

	taskManager, _ := createTestTaskManager(taskRepositoryMock)

	stopErr := taskManager.Stop(0)

	if !errors.Is(stopErr, expectedErr) {
		t.Error("should return an error")
	}
}

func TestTaskManagerStopTimeRangeManagerError(t *testing.T) {
	expectedErr := errors.New("not found")
	taskRepositoryMock := createTaskRepositoryMock()
	manager, timeRangeRepoMock := createTestTaskManager(taskRepositoryMock)

	timeRangeRepoMock.saveError = expectedErr

	task, _ := manager.Create("my task")

	manager.Start(task.GetId())

	if !task.IsRunning() {
		t.Error("the task should be running")
	}

	stopErr := manager.Stop(task.GetId())

	if !errors.Is(stopErr, expectedErr) {
		t.Error("should return an error")
	}
}

func TestTaskManagerList(t *testing.T) {
	taskRepositoryMock := createTaskRepositoryMock()
	manager, _ := createTestTaskManager(taskRepositoryMock)

	taskNames := []string{"one", "two", "three"}

	for _, taskName := range taskNames {
		manager.Create(taskName)
	}

	tasks, listErr := manager.List()

	if listErr != nil {
		t.Error(test.NoError(listErr))
	}

	if len(tasks) != len(taskNames) {
		t.Errorf("should return three values, got %d", len(tasks))
	}

	for _, taskName := range taskNames {
		found := false
		for _, task := range tasks {
			if task.GetName() == entity.ProperNoun(taskName) {
				found = true
			}
		}

		if !found {
			t.Errorf("couldn't find task %s", taskName)
		}
	}
}

func TestTaskManagerListError(t *testing.T) {
	expectedErr := errors.New("listing error")
	taskRepositoryMock := createTaskRepositoryMock()
	taskRepositoryMock.listErr = expectedErr

	manager, _ := createTestTaskManager(taskRepositoryMock)

	if _, listErr := manager.List(); !errors.Is(listErr, expectedErr) {
		t.Error("expected an error")
	}
}
