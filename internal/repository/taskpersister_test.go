package repository

import (
	"errors"
	"testing"

	"github.com/mattn/go-sqlite3"
	"github.com/mehdi-valette/timetracker/internal/entity"
	"github.com/mehdi-valette/timetracker/internal/manager"
	"github.com/mehdi-valette/timetracker/internal/test"
)

func createTestTaskRepo() manager.TaskPersister {
	conn, _ := CreateConnection(":memory:")
	conn.InitializeDb()

	taskRepo := CreateTaskRepository(conn, entity.CreateDate())

	return taskRepo
}

func TestTaskRepositoryCreate(t *testing.T) {
	taskRepo := createTestTaskRepo()

	taskId, createErr := taskRepo.Create()

	if createErr != nil {
		t.Error(test.NoError(createErr))
	}

	if taskId == 0 {
		t.Errorf("task ID should not be 0, got %d", taskId)
	}
}

func TestTaskRepositoryGet(t *testing.T) {
	taskRepo := createTestTaskRepo()

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
	taskRepo := createTestTaskRepo()

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

	taskRepo := createTestTaskRepo()

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

func TestTaskRepositorySaveDuplicate(t *testing.T) {
	name := "my new task"

	taskRepo := createTestTaskRepo()

	taskRepo.Create()
	taskRepo.Create()
	firstTaskId, _ := taskRepo.Create()
	taskRepo.Create()
	secondTaskId, _ := taskRepo.Create()

	firstTask, _ := taskRepo.Get(firstTaskId)
	secondTask, _ := taskRepo.Get(secondTaskId)

	firstTask.Rename(name)
	secondTask.Rename(name)

	if err := taskRepo.Save(firstTask); err != nil {
		t.Error(test.NoError(err))
	}

	err := taskRepo.Save(secondTask)
	var sqlError sqlite3.Error

	if !errors.As(err, &sqlError) || sqlError.ExtendedCode != 2067 {
		t.Error("expected a 'unique constraint' error")
	}
}

func TestTaskRepositorySaveNonExistent(t *testing.T) {
	taskRepo := createTestTaskRepo()

	taskRepo.Create()
	taskRepo.Create()
	taskRepo.Create()

	noTask := entity.CreateTask(-1, "my task", entity.CreateDate())

	if err := taskRepo.Save(noTask); !errors.Is(err, TaskNotFoundErr) {
		t.Error("expected error")
	}
}

func TestTaskRepositoryDelete(t *testing.T) {
	taskRepo := createTestTaskRepo()

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
	taskRepo := createTestTaskRepo()

	taskRepo.Create()
	taskRepo.Create()
	taskRepo.Create()

	if err := taskRepo.Delete(-1); !errors.Is(err, TaskNotFoundErr) {
		t.Error("expected an error")
	}
}

func TestTaskRepositoryList(t *testing.T) {
	taskRepo := createTestTaskRepo()

	names := []string{"task one", "task two", "task three"}

	taskIds := make([]entity.DbId, 0, len(names))
	tasks := make(map[entity.DbId]entity.Tasker, len(names))
	for _, name := range names {
		newId, _ := taskRepo.Create()
		taskIds = append(taskIds, newId)

		newTask := entity.CreateTask(newId, name, entity.CreateDate())
		tasks[newId] = newTask

		taskRepo.Save(newTask)
	}

	taskList, err := taskRepo.List()

	if err != nil {
		t.Error(test.NoError(err))
	}

	if len(taskList) != len(names) {
		t.Errorf("should have 3 tasks, got %d", len(names))
	}

	for _, task := range taskList {
		expectedTask, found := tasks[task.GetId()]

		if !found {
			t.Errorf("should have found the task %d", task.GetId())
		}

		if expectedTask.GetName() != task.GetName() {
			t.Errorf("expected name \"%s\", found \"%s\"", expectedTask.GetName(), task.GetName())
		}
	}
}
