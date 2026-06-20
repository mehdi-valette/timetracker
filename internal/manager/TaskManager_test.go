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
	calls             []mockCalls
	createErr         error
	deleteErr         error
	getByShortNameErr error
	getErr            error
	lastId            entity.DbId
	listErr           error
	saverErr          error
	taskList          map[entity.DbId]entity.Tasker
}

var _ TaskPersister = &taskRepositoryMock{}

// GetByShortName implements [TaskPersister].
func (p *taskRepositoryMock) GetByShortName(shortname string) (entity.Tasker, error) {
	if p.getByShortNameErr != nil {
		return &entity.Task{}, p.getErr
	}

	short := entity.CreateShortName(shortname)

	for _, task := range p.taskList {
		if task.GetShortName() == short {
			mockedTask := entity.CreateTask(task.GetId(), string(task.GetShortName()), string(task.GetName()), entity.CreateDate())
			return mockedTask, nil
		}
	}

	return &entity.Task{}, errors.New("task not found")
}

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

	p.taskList[p.lastId] = entity.CreateTask(p.lastId, "", "", DateMock{})

	return p.lastId, p.createErr
}

func (p *taskRepositoryMock) Save(task entity.Tasker) error {
	p.calls = append(p.calls, mockCalls{"Save", task})

	mockedTask := entity.CreateTask(task.GetId(), string(task.GetShortName()), string(task.GetName()), entity.CreateDate())

	p.taskList[task.GetId()] = mockedTask

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

	mockedTask := entity.CreateTask(task.GetId(), string(task.GetShortName()), string(task.GetName()), entity.CreateDate())

	return mockedTask, nil
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

	task, createError := manager.Create("  Alo  ", "   hello     ")

	if createError != nil {
		t.Error("should not return an error")
	}

	if !taskRepositoryMock.HasBeenCalledWith("Save", task) {
		t.Error("should have called \"save\"")
	}

	if taskRepositoryMock.HasBeenCalledTimes("Save", task) != 1 {
		t.Error("should have called \"save\" 1 time")
	}

	if task.GetShortName() != "alo" {
		t.Errorf("the short name should be \"alo\", got \"%s\"", task.GetShortName())
	}

	if task.GetName() != "hello" {
		t.Errorf("the name should be \"hello\n, it is %s", task.GetName())
	}

	if !taskRepositoryMock.HasBeenCalledWith("Create", nil) {
		t.Error("should have called create")
	}
}

func TestTaskManagerCreateErrorOnCreate(t *testing.T) {
	persistMock := createTaskRepositoryMock()
	expectedErr := errors.New("create error")

	persistMock.createErr = expectedErr

	manager, _ := createTestTaskManager(persistMock)

	task, createError := manager.Create("", "   hello     ")

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

func TestTaskManagerCreateErrorOnSave(t *testing.T) {
	persistMock := createTaskRepositoryMock()
	expectedErr := errors.New("create error")

	persistMock.saverErr = expectedErr

	manager, _ := createTestTaskManager(persistMock)

	task, createError := manager.Create("", "   hello     ")

	if !errors.Is(createError, expectedErr) {
		t.Error("should return an error")
	}

	if task != nil {
		t.Error("task should be nil")
	}

	if !persistMock.HasBeenCalledWith("Create", nil) {
		t.Error("should have called create")
	}

	if list, _ := manager.List(); len(list) != 0 {
		t.Error("should not create a task")
	}
}

func TestTaskManagerCreateErrorOnDelete(t *testing.T) {
	persistMock := createTaskRepositoryMock()
	expectedErr := errors.New("create error")

	persistMock.saverErr = errors.New("")
	persistMock.deleteErr = expectedErr

	manager, _ := createTestTaskManager(persistMock)

	task, createError := manager.Create("", "   hello     ")

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

func TestTaskManagerGetByShortName(t *testing.T) {
	persistMock := createTaskRepositoryMock()

	manager, _ := createTestTaskManager(persistMock)

	createdTask, _ := manager.Create("abc", "hello")

	retrievedTask, getErr := manager.GetByShortName(createdTask.GetShortName().String())

	if getErr != nil {
		t.Error(test.NoError(getErr))
	}

	if retrievedTask.GetId() != createdTask.GetId() || retrievedTask.GetShortName() != createdTask.GetShortName() {
		t.Error("the task recieved is not as expected")
	}
}

func TestTaskManagerGetByShortNameErrorOnTimeRange(t *testing.T) {
	persistMock := createTaskRepositoryMock()

	taskManager, timeRangeManager := createTestTaskManager(persistMock)

	expectedErr := errors.New("time range error")

	timeRangeManager.listByTaskIdError = expectedErr

	taskManager.Create("abc", "hello")

	retrievedTask, getErr := taskManager.GetByShortName("abc")

	if !errors.Is(expectedErr, getErr) {
		t.Error("expected an error")
	}

	if retrievedTask != nil {
		t.Error("the task should be nil")
	}
}

func TestTaskManagerGetByShortNameWithTimerange(t *testing.T) {
	persistMock := createTaskRepositoryMock()

	manager, _ := createTestTaskManager(persistMock)

	createdTask, _ := manager.Create("abc", "hello")

	manager.Start(createdTask.GetId())

	retrievedTask, getByShortNameErr := manager.GetByShortName("abc")

	if getByShortNameErr != nil {
		t.Error(test.NoError(getByShortNameErr))
	}

	if retrievedTask.GetId() != createdTask.GetId() || retrievedTask.GetName() != createdTask.GetName() {
		t.Error("the task recieved is not as expected")
	}

	if !retrievedTask.IsRunning() {
		t.Error("task should be running")
	}
}

func TestTaskManagerGet(t *testing.T) {
	persistMock := createTaskRepositoryMock()

	manager, _ := createTestTaskManager(persistMock)

	createdTask, _ := manager.Create("", "hello")

	retrievedTask, getErr := manager.Get(createdTask.GetId())

	if getErr != nil {
		t.Error(test.NoError(getErr))
	}

	if retrievedTask.GetId() != createdTask.GetId() || retrievedTask.GetName() != createdTask.GetName() {
		t.Error("the task recieved is not as expected")
	}
}

func TestTaskManagerGetErrorOnTimeRange(t *testing.T) {
	persistMock := createTaskRepositoryMock()

	taskManager, timeRangeManager := createTestTaskManager(persistMock)

	expectedErr := errors.New("time range error")

	timeRangeManager.listByTaskIdError = expectedErr

	createdTask, _ := taskManager.Create("", "hello")

	retrievedTask, getErr := taskManager.Get(createdTask.GetId())

	if !errors.Is(expectedErr, getErr) {
		t.Error("expected an error")
	}

	if retrievedTask != nil {
		t.Error("the task should be nil")
	}
}

func TestTaskManagerGetWithTimerange(t *testing.T) {
	persistMock := createTaskRepositoryMock()

	manager, _ := createTestTaskManager(persistMock)

	createdTask, _ := manager.Create("", "hello")

	retrievedTask, startErr := manager.Start(createdTask.GetId())

	if startErr != nil {
		t.Error(test.NoError(startErr))
	}

	if retrievedTask.GetId() != createdTask.GetId() || retrievedTask.GetName() != createdTask.GetName() {
		t.Error("the task recieved is not as expected")
	}

	if !retrievedTask.IsRunning() {
		t.Error("task should be running")
	}
}

func TestTaskManagerGetUnknown(t *testing.T) {
	persistMock := createTaskRepositoryMock()

	manager, _ := createTestTaskManager(persistMock)

	retrievedTask, getErr := manager.Get(-1)

	if getErr == nil {
		t.Error("should return an error")
	}

	if retrievedTask != nil {
		t.Error("task should be empty")
	}
}

func TestTaskManagerStart(t *testing.T) {
	persistMock := createTaskRepositoryMock()

	manager, _ := createTestTaskManager(persistMock)

	task, _ := manager.Create("", "hello")

	if task.IsRunning() {
		t.Error("the task should not yet be running")
	}

	persistMock.calls = []mockCalls{}

	task, _ = manager.Start(task.GetId())

	if !task.IsRunning() {
		t.Error("the task should be running")
	}
}

func TestTaskManagerStartNotFound(t *testing.T) {
	expectedErr := errors.New("task not found")

	persistMock := createTaskRepositoryMock()
	persistMock.getErr = expectedErr

	manager, _ := createTestTaskManager(persistMock)

	task, _ := manager.Create("", "hello")

	if task.IsRunning() {
		t.Error("the task should not yet be running")
	}

	_, startErr := manager.Start(-1)

	if !errors.Is(startErr, expectedErr) {
		t.Error("should return an error")
	}
}

func TestTaskManagerStartErrOnTimeRange(t *testing.T) {
	expectedErr := errors.New("task not found")

	persistMock := createTaskRepositoryMock()

	manager, timeRangeManager := createTestTaskManager(persistMock)
	timeRangeManager.saveError = expectedErr

	task, _ := manager.Create("", "hello")

	if task.IsRunning() {
		t.Error("the task should not yet be running")
	}

	_, startErr := manager.Start(task.GetId())

	if !errors.Is(expectedErr, startErr) {
		t.Error("should return an error")
	}
}

func TestTaskManagerStartAlreadyRunning(t *testing.T) {
	persistMock := createTaskRepositoryMock()

	manager, timeRangeRepository := createTestTaskManager(persistMock)

	task, _ := manager.Create("", "hello")

	if task.IsRunning() {
		t.Error("the task should not yet be running")
	}

	timeRangesBeforeStart, _ := timeRangeRepository.ListByTaskId(task.GetId())

	_, firstStartErr := manager.Start(task.GetId())

	if firstStartErr != nil {
		t.Error("should not return an error")
	}

	timeRangesAfterFirstStart, _ := timeRangeRepository.ListByTaskId(task.GetId())

	_, secondStartErr := manager.Start(task.GetId())

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
		t.Errorf("task should have 1 time range, got: %d", len(timeRangesAfterSecondStart))
	}
}

func TestTaskManagerStartTimeRangeError(t *testing.T) {
	expectedError := errors.New("cannot create")

	persistMock := createTaskRepositoryMock()
	timeRangeRepoMock := createTimeRangeRepositoryMock()
	timeRangeRepoMock.createError = expectedError

	manager := createTestTaskManagerWithMock(persistMock, timeRangeRepoMock)

	task, _ := manager.Create("", "hello")

	_, startErr := manager.Start(task.GetId())

	if !errors.Is(startErr, expectedError) {
		t.Error("should return an error")
	}
}

func TestTaskManagerSave(t *testing.T) {
	persistMock := createTaskRepositoryMock()

	manager, _ := createTestTaskManager(persistMock)

	task, _ := manager.Create("", "my task")

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

	task, _ := manager.Create("", "my task")

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

	task, createErr := manager.Create("", "my task")

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

	task, _ := manager.Create("", "my task")

	task, _ = manager.Start(task.GetId())

	if !task.IsRunning() {
		t.Error("the task should be running")
	}

	_, stopErr := manager.Stop(task.GetId())

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

	task, _ := taskManager.Create("", "my task")

	task, _ = taskManager.Start(task.GetId())

	if !task.IsRunning() {
		t.Error("the task should be running")
	}

	taskManager.Stop(task.GetId())

	if task.IsRunning() {
		t.Error("the task should not be running")
	}

	_, stopErr := taskManager.Stop(task.GetId())

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

	task, _ := taskManager.Create("", "my task")

	if task.IsRunning() {
		t.Error("the task should not be running")
	}

	_, stopErr := taskManager.Stop(task.GetId())

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

	_, stopErr := taskManager.Stop(0)

	if !errors.Is(stopErr, expectedErr) {
		t.Error("should return an error")
	}
}

func TestTaskManagerStopTimeRangeManagerError(t *testing.T) {
	expectedErr := errors.New("not found")
	taskRepositoryMock := createTaskRepositoryMock()
	manager, timeRangeRepoMock := createTestTaskManager(taskRepositoryMock)

	task, _ := manager.Create("", "my task")

	task, _ = manager.Start(task.GetId())

	if !task.IsRunning() {
		t.Error("the task should be running")
	}

	timeRangeRepoMock.saveError = expectedErr

	_, stopErr := manager.Stop(task.GetId())

	if !errors.Is(stopErr, expectedErr) {
		t.Error("should return an error")
	}
}

func TestTaskManagerList(t *testing.T) {
	taskRepositoryMock := createTaskRepositoryMock()
	manager, _ := createTestTaskManager(taskRepositoryMock)

	taskNames := []string{"one", "two", "three"}

	for _, taskName := range taskNames {
		task, _ := manager.Create("", taskName)
		manager.Start(task.GetId())
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
			lastTimeRange, _ := task.GetLastTimeRange()
			if lastTimeRange == nil {
				t.Error("should have a time range")
			}

			if task.GetName() == entity.ProperNoun(taskName) {
				found = true
			}
		}

		if !found {
			t.Errorf("couldn't find task %s", taskName)
		}
	}
}

func TestTaskManagerListErrorOnTimeRangeList(t *testing.T) {
	expectedErr := errors.New("listing error")
	taskRepositoryMock := createTaskRepositoryMock()
	manager, timeRangeMock := createTestTaskManager(taskRepositoryMock)

	timeRangeMock.listError = expectedErr

	if _, listErr := manager.List(); !errors.Is(listErr, expectedErr) {
		t.Error("expected an error")
	}
}

func TestTaskManagerListErrorOnTaskList(t *testing.T) {
	expectedErr := errors.New("listing error")
	taskRepositoryMock := createTaskRepositoryMock()
	taskRepositoryMock.listErr = expectedErr

	manager, _ := createTestTaskManager(taskRepositoryMock)

	if _, listErr := manager.List(); !errors.Is(listErr, expectedErr) {
		t.Error("expected an error")
	}
}

func TestTaskManagerDelete(t *testing.T) {
	taskRepositoryMock := createTaskRepositoryMock()
	manager, _ := createTestTaskManager(taskRepositoryMock)

	taskNames := []string{"one", "two", "three"}

	for _, taskName := range taskNames {
		manager.Create("", taskName)
	}

	tasks, listErr := manager.List()

	if listErr != nil {
		t.Error(test.NoError(listErr))
	}

	if len(tasks) != len(taskNames) {
		t.Errorf("should return three values, got %d", len(tasks))
	}

	if deleteErr := manager.Delete(tasks[1].GetId()); deleteErr != nil {
		t.Error(test.NoError(deleteErr))
	}

	tasksAfter, listErr := manager.List()

	if len(tasksAfter) != len(taskNames)-1 {
		t.Errorf("should have remove a task")
	}
}

func TestTaskManagerDeleteError(t *testing.T) {
	expectedErr := errors.New("error on delete")
	taskRepositoryMock := createTaskRepositoryMock()
	taskRepositoryMock.deleteErr = expectedErr

	manager, _ := createTestTaskManager(taskRepositoryMock)

	taskNames := []string{"one", "two", "three"}

	for _, taskName := range taskNames {
		manager.Create("", taskName)
	}

	tasks, listErr := manager.List()

	if listErr != nil {
		t.Error(test.NoError(listErr))
	}

	if len(tasks) != len(taskNames) {
		t.Errorf("should return three values, got %d", len(tasks))
	}

	if deleteErr := manager.Delete(tasks[1].GetId()); !errors.Is(expectedErr, deleteErr) {
		t.Error(test.NoError(deleteErr))
	}

	tasksAfter, listErr := manager.List()

	if len(tasksAfter) != len(taskNames) {
		t.Errorf("should have remove a task")
	}
}
