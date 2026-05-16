package repository

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/mehdi-valette/timetracker/internal/entity"
	"github.com/mehdi-valette/timetracker/internal/manager"
	"github.com/mehdi-valette/timetracker/internal/test"
)

func createTestTimeRangeRepo() (manager.TimeRangePersister, []entity.Tasker) {
	conn, _ := CreateConnection(":memory:")
	conn.InitializeDb()

	timeRangeRepo := CreateTimeRangeRepository(conn, entity.CreateDate())
	taskRepo := CreateTaskRepository(conn, entity.CreateDate())

	tasks := make([]entity.Tasker, 0, 10)

	for range 10 {
		id, _ := taskRepo.Create()

		name := strconv.FormatUint(rand.Uint64(), 16)
		task := entity.CreateTask(id, name, entity.CreateDate())

		tasks = append(tasks, task)
	}

	return timeRangeRepo, tasks
}

func TestTimeRangeRepositoryCreate(t *testing.T) {
	timeRangeRepo, tasks := createTestTimeRangeRepo()

	timeRangeId, createErr := timeRangeRepo.Create(tasks[0].GetId())

	if createErr != nil {
		t.Error(test.NoError(createErr))
	}

	if timeRangeId == 0 {
		t.Error("ID should not be 0")
	}
}

func TestTimeRangeRepositoryGet(t *testing.T) {
	timeRangeRepo, tasks := createTestTimeRangeRepo()

	timeRangeId, _ := timeRangeRepo.Create(tasks[0].GetId())

	timeRange, getErr := timeRangeRepo.Get(timeRangeId)

	if getErr != nil {
		t.Error(test.NoError(getErr))
	}

	if timeRange.GetId() != timeRangeId {
		t.Errorf("expected ID %d, got %d", timeRangeId, timeRange.GetId())
	}
}

func TestTimeRangeRepositorySave(t *testing.T) {
	timeRangeRepo, tasks := createTestTimeRangeRepo()

	timeRangeId, _ := timeRangeRepo.Create(tasks[0].GetId())

	timeRange := entity.CreateTimeRange(timeRangeId, tasks[0].GetId(), entity.CreateDate())

	if err := timeRangeRepo.Save(timeRange); err != nil {
		t.Error(test.NoError(err))
	}

	panic("verify that the saved value is correct, with and without start and end dates")
}
