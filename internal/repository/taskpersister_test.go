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

func TestTaskRepositoryGetByShortName(t *testing.T) {
	taskRepo := createTestTaskRepo()
	shortname := "TestTaskRepositoryGetByShortName"

	taskRepo.Create()
	taskRepo.Create()
	taskId, _ := taskRepo.Create()
	taskRepo.Create()

	task, _ := taskRepo.Get(taskId)

	task.Reshortname(shortname)
	taskRepo.Save(task)

	task, getErr := taskRepo.GetByShortName(shortname)

	if getErr != nil {
		t.Error(test.NoError(getErr))
	}

	if task.GetShortName() != entity.CreateShortName(shortname) {
		t.Errorf("expected short-name %s, got %s", task.GetShortName(), shortname)
	}
}

func TestTaskRepositoryGetByShortNameNotFound(t *testing.T) {
	taskRepo := createTestTaskRepo()
	shortname := "TestTaskRepositoryGetByShortNameNotFound"

	taskRepo.Create()
	taskRepo.Create()
	taskRepo.Create()

	if _, getErr := taskRepo.GetByShortName(shortname); !errors.Is(getErr, manager.TaskNotFoundErr) {
		t.Error("should return an error")
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

	if _, getErr := taskRepo.Get(-1); !errors.Is(getErr, manager.TaskNotFoundErr) {
		t.Error("should return an error")
	}
}

func TestTaskRepositorySave(t *testing.T) {
	firstExpectedShortName := "abc"
	firstExpectedName := "my new task"
	secondExpectedShortName := "def"
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
	firstTask.Reshortname(firstExpectedShortName)
	secondTask.Rename(secondExpectedName)
	secondTask.Reshortname(secondExpectedShortName)

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

	if firstTaskAfterSave.GetShortName() != entity.ShortName(firstExpectedShortName) {
		t.Errorf("should have the short name \"%s\", got \"%s\"", firstExpectedShortName, firstTaskAfterSave.GetShortName())
	}

	if secondTaskAfterSave.GetName() != entity.ProperNoun(secondExpectedName) {
		t.Errorf("should have the name \"%s\", got \"%s\"", secondExpectedName, secondTaskAfterSave.GetName())
	}

	if secondTaskAfterSave.GetShortName() != entity.ShortName(secondExpectedShortName) {
		t.Errorf("should have the short name \"%s\", got \"%s\"", secondExpectedShortName, secondTaskAfterSave.GetShortName())
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

	noTask := entity.CreateTask(-1, "", "my task", entity.CreateDate())

	if err := taskRepo.Save(noTask); !errors.Is(err, manager.TaskNotFoundErr) {
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

	if err := taskRepo.Delete(-1); !errors.Is(err, manager.TaskNotFoundErr) {
		t.Error("expected an error")
	}
}

func TestTaskRepositoryList(t *testing.T) {
	taskRepo := createTestTaskRepo()

	names := []struct {
		short string
		name  string
	}{{"one", "task one"}, {"two", "task two"}, {"three", "task three"}}

	taskIds := make([]entity.DbId, 0, len(names))
	tasks := make(map[entity.DbId]entity.Tasker, len(names))
	for _, name := range names {
		newId, _ := taskRepo.Create()
		taskIds = append(taskIds, newId)

		newTask := entity.CreateTask(newId, name.short, name.name, entity.CreateDate())
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

		if expectedTask.GetShortName() != task.GetShortName() {
			t.Errorf("expected short name \"%s\", found \"%s\"", expectedTask.GetShortName(), task.GetShortName())
		}
	}
}
