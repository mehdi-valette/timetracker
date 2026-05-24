package repository

import (
	"errors"
	"math/rand"
	"strconv"
	"testing"
	"time"

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

func TestTimeRangeRepositorySaveAndGet(t *testing.T) {
	// prepare the repository
	timeRangeRepo, tasks := createTestTimeRangeRepo()
	timeRangeId, _ := timeRangeRepo.Create(tasks[0].GetId())

	// prepare the dates
	startDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.Add(time.Second * 30)
	date := entity.CreateDateMock(startDate)

	// create the time range
	timeRange := entity.CreateTimeRange(timeRangeId, tasks[0].GetId(), &date)

	if err := timeRangeRepo.Save(timeRange); err != nil {
		t.Error(test.NoError(err))
	}

	// modify and save the time range several times
	timeRangeBeforeStart, _ := timeRangeRepo.Get(timeRange.GetId())

	timeRange.Start()
	timeRangeRepo.Save(timeRange)

	timeRangeAfterStart, _ := timeRangeRepo.Get(timeRange.GetId())

	date.Set(endDate)
	timeRange.End()
	timeRangeRepo.Save(timeRange)

	timeRangeAfterEnd, _ := timeRangeRepo.Get(timeRange.GetId())

	// verify the results
	if timeRangeBeforeStart.HasStarted() {
		t.Error("should not have started")
	}
	if timeRangeBeforeStart.HasEnded() {
		t.Error("should not have ended")
	}

	if !timeRangeAfterStart.HasStarted() {
		t.Error("should have started")
	}
	if timeRangeAfterStart.HasEnded() {
		t.Error("should not have ended")
	}
	if timeRangeAfterStart.GetStart().GetSeconds() != startDate.Unix() {
		t.Error("should have started on time")
	}

	if !timeRangeAfterEnd.HasStarted() {
		t.Error("should have started")
	}
	if !timeRangeAfterEnd.HasEnded() {
		t.Error("should have ended")
	}
	if timeRangeAfterEnd.GetStart().GetSeconds() != startDate.Unix() {
		t.Error("should have started on time")
	}
	if timeRangeAfterEnd.GetEnd().GetSeconds() != endDate.Unix() {
		t.Error("should have ended on time")
	}
}

func TestTimeRangeRepositoryGetNotFound(t *testing.T) {
	timeRangeRepo, _ := createTestTimeRangeRepo()

	_, getErr := timeRangeRepo.Get(-1)

	if !errors.Is(getErr, TimeRangeNotFoundErr) {
		t.Error("Expected an error")
	}
}

func TestTimeRangeRepositoryListByTaskId(t *testing.T) {
	timeRangeRepo, tasks := createTestTimeRangeRepo()

	date := entity.CreateDateMock(time.Now())

	expectedCount := 4
	expectedRanges := make([]entity.TimeRanger, 0, expectedCount)

	for range expectedCount {
		timeRangeId, _ := timeRangeRepo.Create(tasks[0].GetId())
		timeRange := entity.CreateTimeRange(timeRangeId, tasks[0].GetId(), &date)

		date.Set(time.Unix(rand.Int63(), 0))
		timeRange.Start()

		timeRangeRepo.Save(timeRange)

		expectedRanges = append(expectedRanges, timeRange)
	}

	timeRanges, listErr := timeRangeRepo.ListByTaskId(tasks[0].GetId())

	if listErr != nil {
		t.Error(test.NoError(listErr))
	}

	if len(timeRanges) != expectedCount {
		t.Errorf("got %d errors, expected %d", len(timeRanges), expectedCount)
	}

	for _, timeRange := range timeRanges {
		found := false
		for _, expectedRange := range expectedRanges {
			if expectedRange.GetId() == timeRange.GetId() {
				found = true

				if expectedRange.GetStart().GetSeconds() != timeRange.GetStart().GetSeconds() {
					t.Errorf("expect time range %d to start at %d, but started at %d",
						expectedRange.GetId(),
						expectedRange.GetStart().GetSeconds(),
						timeRange.GetStart().GetSeconds(),
					)
				}
			}
		}

		if !found {
			t.Errorf("didn't expect time range with ID %d", timeRange.GetId())
		}
	}
}

func TestTimeRangeRepositoryListByTaskIdEmpty(t *testing.T) {
	timeRangeRepo, tasks := createTestTimeRangeRepo()

	timeRanges, listErr := timeRangeRepo.ListByTaskId(tasks[0].GetId())

	if listErr != nil {
		t.Error(test.NoError(listErr))
	}

	if len(timeRanges) != 0 {
		t.Errorf("expected an empty list, got %d", len(timeRanges))
	}
}

func TestTimeRangeRepositoryListByTaskIdNotFound(t *testing.T) {
	timeRangeRepo, _ := createTestTimeRangeRepo()

	timeRanges, listErr := timeRangeRepo.ListByTaskId(-1)

	if listErr != nil {
		t.Error(test.NoError(listErr))
	}

	if len(timeRanges) != 0 {
		t.Errorf("expected an empty list, got %d", len(timeRanges))
	}
}
