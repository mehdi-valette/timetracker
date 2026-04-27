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

func TestTaskManagerCreate(t *testing.T) {
	persistMock := createPersisterMock()

	manager := CreateTaskManager(persistMock, entity.CreateDate())

	task, createError := manager.Create("   hello     ")

	if createError != nil {
		t.Error("should not return an error")
	}

	if task.Name() != "hello" {
		t.Errorf("the name should be \"hello\n, it is %s", task.Name())
	}

	if !persistMock.HasBeenCalledWith("Create", nil) {
		t.Error("should have called create")
	}
}

func TestTaskManagerCreateError(t *testing.T) {
	persistMock := createPersisterMock()
	expectedErr := errors.New("create error")

	persistMock.createErr = expectedErr

	manager := CreateTaskManager(persistMock, entity.CreateDate())

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
