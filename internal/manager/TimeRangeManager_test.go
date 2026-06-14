package manager

import (
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/mehdi-valette/timetracker/internal/entity"
	"github.com/mehdi-valette/timetracker/internal/test"
)

type DateMock struct {
	currentTime time.Time
}

var _ entity.Dater = DateMock{}

func (d DateMock) Now() entity.Timestamp {
	return entity.Timestamp(d.currentTime)
}

func (d *DateMock) Set(t time.Time) {
	d.currentTime = t
}

type MethodCall struct {
	method string
	param  any
}

type timeRangeRepositoryMock struct {
	calledWith          []MethodCall
	countCreateCalls    uint8
	countDeleteCalls    uint8
	countSaveCalls      uint8
	createError         error
	deleteError         error
	getErr              error
	saveError           error
	listByTaskIdError   error
	listError           error
	getLastTimeRangeErr error
	lastId              entity.DbId
	timeRanges          map[entity.DbId]entity.TimeRanger
}

// GetOpenTimeRanges implements [TimeRangePersister].
func (trrm *timeRangeRepositoryMock) GetLastTimeRange() (entity.TimeRanger, error) {
	if trrm.getLastTimeRangeErr != nil {
		return nil, trrm.getLastTimeRangeErr
	}

	lastEnd := 0
	var lastTimeRange entity.TimeRanger = nil

	for _, timeRange := range trrm.timeRanges {
		if timeRange.HasStarted() && !timeRange.HasEnded() {
			return timeRange, nil
		}

		if int(timeRange.GetEnd().GetSeconds()) > lastEnd {
			lastEnd = int(timeRange.GetEnd().GetSeconds())
			lastTimeRange = timeRange
		}
	}

	return lastTimeRange, nil
}

func createTimeRangeRepositoryMock() *timeRangeRepositoryMock {
	return &timeRangeRepositoryMock{
		timeRanges: make(map[entity.DbId]entity.TimeRanger),
	}
}

func (trrm *timeRangeRepositoryMock) Create(taskId entity.DbId) (entity.DbId, error) {
	trrm.lastId += 1
	trrm.countCreateCalls += 1

	trrm.calledWith = append(trrm.calledWith, MethodCall{"Create", nil})

	timeRange := entity.CreateTimeRange(trrm.lastId, taskId, DateMock{})
	trrm.timeRanges[timeRange.GetId()] = timeRange

	return trrm.lastId, trrm.createError
}

func (trrm *timeRangeRepositoryMock) Save(timeRange entity.TimeRanger) error {
	trrm.countSaveCalls += 1

	trrm.calledWith = append(trrm.calledWith, MethodCall{"Save", timeRange})

	trrm.timeRanges[timeRange.GetId()] = timeRange

	return trrm.saveError
}

func (trrm *timeRangeRepositoryMock) Get(timeRangeId entity.DbId) (entity.TimeRanger, error) {
	if trrm.getErr != nil {
		return &entity.TimeRange{}, trrm.getErr
	}

	timeRange, found := trrm.timeRanges[timeRangeId]

	if !found {
		return &entity.TimeRange{}, errors.New("not found")
	}

	return timeRange, nil
}

func (trrm *timeRangeRepositoryMock) Delete(id entity.DbId) error {
	trrm.countDeleteCalls += 1

	trrm.calledWith = append(trrm.calledWith, MethodCall{"Delete", id})

	return trrm.deleteError
}

func (trrm *timeRangeRepositoryMock) HasBeenCalledWith(method string, param any) bool {
	for _, call := range trrm.calledWith {
		if call.method == method && call.param == param {
			return true
		}
	}

	return false
}

func (trrm *timeRangeRepositoryMock) ListByTaskId(taskId entity.DbId) ([]entity.TimeRanger, error) {
	if trrm.listByTaskIdError != nil {
		return []entity.TimeRanger{}, trrm.listByTaskIdError
	}

	timeRanges := make([]entity.TimeRanger, 0, len(trrm.timeRanges))

	for _, timeRange := range trrm.timeRanges {
		if timeRange.GetTaskId() == taskId {
			timeRanges = append(timeRanges, timeRange)
		}
	}

	return timeRanges, nil
}

func (trrm *timeRangeRepositoryMock) List() ([]entity.TimeRanger, error) {
	if trrm.listError != nil {
		return []entity.TimeRanger{}, trrm.listError
	}

	timeRanges := make([]entity.TimeRanger, 0, len(trrm.timeRanges))

	for _, timeRange := range trrm.timeRanges {
		timeRanges = append(timeRanges, timeRange)
	}

	return timeRanges, nil
}

var _ TimeRangePersister = &timeRangeRepositoryMock{}

func TestTimeRangeManagerCreate(t *testing.T) {
	repoMock := createTimeRangeRepositoryMock()

	manager := CreateTimeRangeManager(repoMock, DateMock{})

	timeRange, createErr := manager.Create(0)

	if createErr != nil {
		t.Errorf("should not have returned an error")
	}

	if timeRange.GetId() != repoMock.lastId {
		t.Errorf("should have returned the time range of the repo")
	}

	if timeRange.HasStarted() {
		t.Errorf("time range should not have started")
	}

	if repoMock.countCreateCalls != 1 {
		t.Errorf("should have called create from the repository")
	}

	if repoMock.countSaveCalls != 1 {
		t.Errorf("should have called save from the repository")
	}
}

func TestTimeRangeManagerCreateError(t *testing.T) {
	chosenError := errors.New("create error")

	repoMock := createTimeRangeRepositoryMock()
	repoMock.createError = chosenError

	manager := CreateTimeRangeManager(repoMock, DateMock{})

	timeRange, createErr := manager.Create(0)

	if !errors.Is(createErr, chosenError) {
		t.Error("should have returned an error")
	}

	if timeRange != nil {
		t.Error("should not have created a time range")
	}
}

func TestTimeRangeManagerStart(t *testing.T) {
	repoMock := createTimeRangeRepositoryMock()

	manager := CreateTimeRangeManager(repoMock, DateMock{})

	timeRange, createErr := manager.Create(0)

	if createErr != nil {
		t.Error("should not return an error")
	}

	_, startErr := manager.Start(timeRange.GetId())

	if startErr != nil {
		t.Error("should not return an error")
	}

	if !timeRange.HasStarted() {
		t.Error("time range should have started")
	}

	if !repoMock.HasBeenCalledWith("Save", timeRange) {
		t.Error("should have saved the time range")
	}

}

func TestTimeRangeManagerStartNotFound(t *testing.T) {
	repoMock := createTimeRangeRepositoryMock()

	manager := CreateTimeRangeManager(repoMock, DateMock{})

	_, createErr := manager.Create(0)

	if createErr != nil {
		t.Error("should not return an error")
	}

	_, startErr := manager.Start(-1)

	if startErr == nil {
		t.Error("should return an error")
	}
}

func TestTimeRangeManagerStartError(t *testing.T) {
	chosenError := errors.New("save error")

	repoMock := createTimeRangeRepositoryMock()
	repoMock.saveError = chosenError

	manager := CreateTimeRangeManager(repoMock, DateMock{})

	timeRange, createErr := manager.Create(0)

	if createErr != nil {
		t.Error("should not return an error")
	}

	_, startErr := manager.Start(timeRange.GetId())

	if !errors.Is(startErr, chosenError) {
		t.Error("should have returned an error")
	}

	if !timeRange.HasStarted() {
		t.Error("time range should not have started")
	}
}

func TestTimeRangeManagerStop(t *testing.T) {
	repoMock := createTimeRangeRepositoryMock()

	manager := CreateTimeRangeManager(repoMock, DateMock{})

	timeRange, createErr := manager.Create(0)

	if createErr != nil {
		t.Error("should not return an error")
	}

	if timeRange.HasEnded() {
		t.Error("time range should not have ended")
	}

	stopErr := manager.Stop(timeRange.GetId())

	if stopErr != nil {
		t.Error("should not return an error")
	}

	if !timeRange.HasEnded() {
		t.Error("time range should have ended")
	}

	if !repoMock.HasBeenCalledWith("Save", timeRange) {
		t.Error("should have saved the time range")
	}
}

func TestTimeRangeManagerStopError(t *testing.T) {
	chosenError := errors.New("save error")

	repoMock := createTimeRangeRepositoryMock()
	repoMock.saveError = chosenError

	manager := CreateTimeRangeManager(repoMock, DateMock{})

	timeRange, createErr := manager.Create(0)

	if createErr != nil {
		t.Error("should not return an error")
	}

	if timeRange.HasEnded() {
		t.Error("time range should not have ended")
	}

	stopErr := manager.Stop(timeRange.GetId())

	if !errors.Is(stopErr, chosenError) {
		t.Error("should have returned an error")
	}

	if !timeRange.HasEnded() {
		t.Error("time range should have ended")
	}
}

func TestTimeRangeManagerStopNotFound(t *testing.T) {
	repoMock := createTimeRangeRepositoryMock()

	manager := CreateTimeRangeManager(repoMock, DateMock{})

	_, createErr := manager.Create(0)

	if createErr != nil {
		t.Errorf("should not return an error")
	}

	stopErr := manager.Stop(-1)

	if stopErr == nil {
		t.Errorf("should return an error")
	}
}

func TestTimeRangeManagerSave(t *testing.T) {
	repoMock := createTimeRangeRepositoryMock()

	manager := CreateTimeRangeManager(repoMock, DateMock{}).(*TimeRangeManagement)

	timeRange, createErr := manager.Create(0)

	if createErr != nil {
		t.Error("should not return an error")
	}

	errorSave := manager.Save(timeRange)

	if errorSave != nil {
		t.Error("should not have returned an error")
	}

	if !repoMock.HasBeenCalledWith("Save", timeRange) {
		t.Error("should have called save")
	}
}

func TestTimeRangeManagerSaveError(t *testing.T) {
	expectedErr := errors.New("delete error")

	repoMock := createTimeRangeRepositoryMock()
	repoMock.saveError = expectedErr

	manager := CreateTimeRangeManager(repoMock, DateMock{}).(*TimeRangeManagement)

	timeRange, createErr := manager.Create(0)

	if createErr != nil {
		t.Error("should not return an error")
	}

	errorDelete := manager.Save(timeRange)

	if !errors.Is(errorDelete, expectedErr) {
		t.Error("should have returned a delete error")
	}
}

func TestTimeRangeManagerList(t *testing.T) {
	date := entity.CreateDateMock(time.Now())
	repoMock := createTimeRangeRepositoryMock()

	manager := CreateTimeRangeManager(repoMock, &date).(*TimeRangeManagement)

	expectedCount := 100
	expectedTimeRanges := make([]entity.TimeRanger, 0, expectedCount)

	for range expectedCount {
		timeRange, _ := manager.Create(entity.DbId(rand.Int31n(10)))
		date.Set(time.Unix(rand.Int63(), 0))

		timeRange.Start()
		manager.Save(timeRange)

		expectedTimeRanges = append(expectedTimeRanges, timeRange)
	}

	listedTimeRanges, listErr := manager.List()

	if listErr != nil {
		t.Error(test.NoError(listErr))
	}

	if len(listedTimeRanges) != expectedCount {
		t.Errorf("expected %d records, got %d", expectedCount, len(listedTimeRanges))
	}

	for _, listedTimeRange := range listedTimeRanges {
		found := false

		for _, expectedTimeRange := range expectedTimeRanges {
			if expectedTimeRange.GetId() == listedTimeRange.GetId() {
				found = true

				if expectedTimeRange.GetStart().GetSeconds() != listedTimeRange.GetStart().GetSeconds() {
					t.Errorf("the time range %d should start at %d, but got %d",
						listedTimeRange.GetId(),
						expectedTimeRange.GetStart().GetSeconds(),
						listedTimeRange.GetStart().GetSeconds(),
					)
				}
			}
		}

		if !found {
			t.Errorf("cannot find time range %d", listedTimeRange.GetId())
		}
	}
}

func TestTimeRangeManagerListByTaskId(t *testing.T) {
	date := entity.CreateDateMock(time.Now())
	repoMock := createTimeRangeRepositoryMock()

	manager := CreateTimeRangeManager(repoMock, &date).(*TimeRangeManagement)

	expectedCount := 10
	expectedTimeRanges := make([]entity.TimeRanger, 0, expectedCount)

	for range expectedCount {
		timeRange, _ := manager.Create(0)
		date.Set(time.Unix(rand.Int63(), 0))

		timeRange.Start()
		manager.Save(timeRange)

		expectedTimeRanges = append(expectedTimeRanges, timeRange)
	}

	listedTimeRanges, listErr := manager.ListByTaskId(0)

	if listErr != nil {
		t.Error(test.NoError(listErr))
	}

	if len(listedTimeRanges) != expectedCount {
		t.Errorf("expected %d records, got %d", expectedCount, len(listedTimeRanges))
	}

	for _, listedTimeRange := range listedTimeRanges {
		found := false

		for _, expectedTimeRange := range expectedTimeRanges {
			if expectedTimeRange.GetId() == listedTimeRange.GetId() {
				found = true

				if expectedTimeRange.GetStart().GetSeconds() != listedTimeRange.GetStart().GetSeconds() {
					t.Errorf("the time range %d should start at %d, but got %d",
						listedTimeRange.GetId(),
						expectedTimeRange.GetStart().GetSeconds(),
						listedTimeRange.GetStart().GetSeconds(),
					)
				}
			}
		}

		if !found {
			t.Errorf("cannot find time range %d", listedTimeRange.GetId())
		}
	}
}

func TestTimeRangeManagerListByTaskIdError(t *testing.T) {
	expectedErr := errors.New("listByTaskId error")

	date := entity.CreateDateMock(time.Now())
	repoMock := createTimeRangeRepositoryMock()

	repoMock.listByTaskIdError = expectedErr

	manager := CreateTimeRangeManager(repoMock, &date).(*TimeRangeManagement)

	listedTimeRanges, listErr := manager.ListByTaskId(0)

	if !errors.Is(listErr, expectedErr) {
		t.Error("should return the expected error")
	}

	if len(listedTimeRanges) != 0 {
		t.Error("should return an empty list")
	}
}

func TestTimeRangeManagerGetOpenTimeRanges(t *testing.T) {
	date := entity.CreateDateMock(time.Now())
	repoMock := createTimeRangeRepositoryMock()

	manager := CreateTimeRangeManager(repoMock, &date).(*TimeRangeManagement)

	expectedCount := 10
	expectedTimeRanges := make([]entity.TimeRanger, 0, expectedCount)

	for index := range expectedCount {
		timeRange, _ := manager.Create(0)
		date.Set(time.Unix(rand.Int63n(1000), 0))

		timeRange.Start()

		if index != 0 {
			date.Set(time.Unix(rand.Int63n(1000000000), 0))
			timeRange.End()
		}

		manager.Save(timeRange)

		expectedTimeRanges = append(expectedTimeRanges, timeRange)
	}

	lastRange, lastRangeErr := manager.GetLastTimeRange()

	if lastRangeErr != nil {
		t.Error(test.NoError(lastRangeErr))
	}

	if lastRange.GetId() != 1 {
		t.Errorf("expected time range 1, got %d", lastRange.GetId())
	}
}

func TestTimeRangeManagerGetOpenTimeRangesError(t *testing.T) {
	expectedErr := errors.New("open range error")
	date := entity.CreateDateMock(time.Now())
	repoMock := createTimeRangeRepositoryMock()

	manager := CreateTimeRangeManager(repoMock, &date).(*TimeRangeManagement)
	repoMock.getLastTimeRangeErr = expectedErr

	lastRange, lastRangeErr := manager.GetLastTimeRange()

	if !errors.Is(lastRangeErr, expectedErr) {
		t.Error("expected an error")
	}

	if lastRange != nil {
		t.Errorf("expected nil, got %+v", lastRange)
	}
}
