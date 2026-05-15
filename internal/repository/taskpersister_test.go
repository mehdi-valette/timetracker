package repository

import (
	"errors"
	"testing"

	"github.com/mehdi-valette/timetracker/internal/entity"
	"github.com/mehdi-valette/timetracker/internal/test"
)

func TestTaskRepositoryCreate(t *testing.T) {
	conn := DbConnection{}
	conn.Connect(":memory:")
	conn.InitializeDb()

	taskRepo := CreateTaskRepository(&conn, entity.CreateDate())

	taskId, createErr := taskRepo.Create()

	if createErr != nil {
		t.Error(test.NoError(createErr))
	}

	if taskId == 0 {
		t.Errorf("task ID should not be 0, got %d", taskId)
	}
}

func TestTaskRepositoryGet(t *testing.T) {
	conn := DbConnection{}
	conn.Connect(":memory:")
	conn.InitializeDb()

	taskRepo := CreateTaskRepository(&conn, entity.CreateDate())

	taskRepo.Create()
	taskRepo.Create()
	taskId, _ := taskRepo.Create()
	taskRepo.Create()

	task, getErr := taskRepo.Get(taskId)

	if getErr != nil {
		t.Error(test.NoError(getErr))
	}

	if task.GetId() != taskId {
		t.Errorf("got ID %d, should be %d", task.GetId(), taskId)
	}
}

func TestTaskRepositoryGetNonExistent(t *testing.T) {
	conn := DbConnection{}
	conn.Connect(":memory:")
	conn.InitializeDb()

	taskRepo := CreateTaskRepository(&conn, entity.CreateDate())

	taskRepo.Create()
	taskRepo.Create()
	taskRepo.Create()

	if _, getErr := taskRepo.Get(-1); !errors.Is(getErr, TaskNotFoundErr) {
		t.Error("should return an error")
	}
}

func TestTaskRepositorySave(t *testing.T) {
	firstExpectedName := "my new task"
	secondExpectedName := "my second new task"

	conn := DbConnection{}
	conn.Connect(":memory:")
	conn.InitializeDb()

	taskRepo := CreateTaskRepository(&conn, entity.CreateDate())

	taskRepo.Create()
	taskRepo.Create()
	firstTaskId, _ := taskRepo.Create()
	taskRepo.Create()
	secondTaskId, _ := taskRepo.Create()

	firstTask, _ := taskRepo.Get(firstTaskId)
	secondTask, _ := taskRepo.Get(secondTaskId)

	firstTask.Rename(firstExpectedName)
	secondTask.Rename(secondExpectedName)

	if err := taskRepo.Save(firstTask); err != nil {
		t.Error(test.NoError(err))
	}

	if err := taskRepo.Save(secondTask); err != nil {
		t.Error(test.NoError(err))
	}

	firstTaskAfterSave, _ := taskRepo.Get(firstTaskId)
	secondTaskAfterSave, _ := taskRepo.Get(secondTaskId)

	if firstTaskAfterSave.GetName() != entity.ProperNoun(firstExpectedName) {
		t.Errorf("should have the name \"%s\", got \"%s\"", firstExpectedName, firstTaskAfterSave.GetName())
	}

	if secondTaskAfterSave.GetName() != entity.ProperNoun(secondExpectedName) {
		t.Errorf("should have the name \"%s\", got \"%s\"", secondExpectedName, secondTaskAfterSave.GetName())
	}
}

func TestTaskRepositorySaveNonExistent(t *testing.T) {
	conn := DbConnection{}
	conn.Connect(":memory:")
	conn.InitializeDb()

	taskRepo := CreateTaskRepository(&conn, entity.CreateDate())

	taskRepo.Create()
	taskRepo.Create()
	taskRepo.Create()

	noTask := entity.CreateTask(-1, "my task", entity.CreateDate())

	if err := taskRepo.Save(noTask); !errors.Is(err, TaskNotFoundErr) {
		t.Error("expected error")
	}
}

func TestTaskRepositoryDelete(t *testing.T) {
	conn := DbConnection{}
	conn.Connect(":memory:")
	conn.InitializeDb()

	taskRepo := CreateTaskRepository(&conn, entity.CreateDate())

	taskRepo.Create()
	taskRepo.Create()
	fristTaskId, _ := taskRepo.Create()
	taskRepo.Create()
	secondTaskId, _ := taskRepo.Create()
	taskRepo.Create()

	// verify the tasks exist
	if _, err := taskRepo.Get(fristTaskId); err != nil {
		t.Error(test.NoError(err))
	}

	if _, err := taskRepo.Get(secondTaskId); err != nil {
		t.Error(test.NoError(err))
	}

	// delete one task
	if err := taskRepo.Delete(fristTaskId); err != nil {
		t.Error(test.NoError(err))
	}

	// verify the deleted task is gone and the other still exists
	if _, err := taskRepo.Get(fristTaskId); err == nil {
		t.Error("expected an error")
	}

	if _, err := taskRepo.Get(secondTaskId); err != nil {
		t.Error(test.NoError(err))
	}
}

func TestTaskRepositoryDeleteNonExistent(t *testing.T) {
	conn := DbConnection{}
	conn.Connect(":memory:")
	conn.InitializeDb()

	taskRepo := CreateTaskRepository(&conn, entity.CreateDate())

	taskRepo.Create()
	taskRepo.Create()
	taskRepo.Create()

	if err := taskRepo.Delete(-1); !errors.Is(err, TaskNotFoundErr) {
		t.Error("expected an error")
	}
}
