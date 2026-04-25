package manager

import (
	"errors"
	"testing"
	"time"

	"github.com/mehdi-valette/timetracker/internal/entity"
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

type TimeRangeRepositoryMock struct {
	lastId           entity.DbId
	deleteError      error
	countCreateCalls uint8
	countSaveCalls   uint8
	countDeleteCalls uint8
	calledWith       []MethodCall
}

var _ TimeRangeRepositoryManager = &TimeRangeRepositoryMock{}

func (trrm *TimeRangeRepositoryMock) Create() (entity.DbId, error) {
	trrm.lastId += 1
	trrm.countCreateCalls += 1

	trrm.calledWith = append(trrm.calledWith, MethodCall{"Create", nil})

	return trrm.lastId, nil
}

func (trrm *TimeRangeRepositoryMock) Save(timeRange entity.TimeRanger) error {
	trrm.countSaveCalls += 1

	trrm.calledWith = append(trrm.calledWith, MethodCall{"Save", timeRange})

	return nil
}

func (trrm *TimeRangeRepositoryMock) Delete(id entity.DbId) error {
	trrm.countDeleteCalls += 1

	trrm.calledWith = append(trrm.calledWith, MethodCall{"Delete", id})

	return trrm.deleteError
}

func (trrm *TimeRangeRepositoryMock) HasBeenCalledWith(method string, param any) bool {
	for _, call := range trrm.calledWith {
		if call.method == method && call.param == param {
			return true
		}
	}

	return false
}

func TestTimeRangeManagerCreate(t *testing.T) {
	repoMock := &TimeRangeRepositoryMock{}

	manager := CreateTimeRangeManager(repoMock, DateMock{})

	timeRange, createErr := manager.Create()

	if createErr != nil {
		t.Errorf("should not have returned an error")
	}

	if timeRange.GetId() != repoMock.lastId {
		t.Errorf("should have returned the time range of the repo")
	}

	if timeRange.HasStated() {
		t.Errorf("time range should not have started")
	}

	if repoMock.countCreateCalls != 1 {
		t.Errorf("should have called create from the repository")
	}

	if repoMock.countSaveCalls != 1 {
		t.Errorf("should have called save from the repository")
	}
}

func TestTimeRangeManagerStart(t *testing.T) {
	repoMock := &TimeRangeRepositoryMock{}

	manager := CreateTimeRangeManager(repoMock, DateMock{})

	timeRange, createErr := manager.Create()

	if createErr != nil {
		t.Error("should not return an error")
	}

	startErr := manager.Start(timeRange.GetId())

	if startErr != nil {
		t.Error("should not return an error")
	}

	if !timeRange.HasStated() {
		t.Error("time range should have started")
	}

	if !repoMock.HasBeenCalledWith("Save", timeRange) {
		t.Error("should have saved the time range")
	}

}

func TestTimeRangeManagerStartNotFound(t *testing.T) {
	repoMock := &TimeRangeRepositoryMock{}

	manager := CreateTimeRangeManager(repoMock, DateMock{})

	_, createErr := manager.Create()

	if createErr != nil {
		t.Error("should not return an error")
	}

	startErr := manager.Start(-1)

	if startErr == nil {
		t.Error("should return an error")
	}
}

func TestTimeRangeManagerStop(t *testing.T) {
	repoMock := &TimeRangeRepositoryMock{}

	manager := CreateTimeRangeManager(repoMock, DateMock{})

	timeRange, createErr := manager.Create()

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

func TestTimeRangeManagerStopNotFound(t *testing.T) {
	repoMock := &TimeRangeRepositoryMock{}

	manager := CreateTimeRangeManager(repoMock, DateMock{})

	_, createErr := manager.Create()

	if createErr != nil {
		t.Errorf("should not return an error")
	}

	stopErr := manager.Stop(-1)

	if stopErr == nil {
		t.Errorf("should return an error")
	}
}

func TestTimeRangeManagerDelete(t *testing.T) {
	repoMock := &TimeRangeRepositoryMock{}

	manager := CreateTimeRangeManager(repoMock, DateMock{}).(*TimeRangeManagement)

	timeRange, createErr := manager.Create()

	if createErr != nil {
		t.Error("should not return an error")
	}

	_, foundBeforeDelete := manager.timeRanges[timeRange.GetId()]

	if !foundBeforeDelete {
		t.Error("should find the time range")
	}

	if repoMock.countDeleteCalls != 0 {
		t.Error("should not have called delete yet")
	}

	deleteError := manager.Delete(timeRange.GetId())

	if deleteError != nil {
		t.Error("should not have returned an error")
	}

	_, foundAfterDelete := manager.timeRanges[timeRange.GetId()]

	if foundAfterDelete {
		t.Error("should not find the time range")
	}

	if !repoMock.HasBeenCalledWith("Delete", timeRange.GetId()) {
		t.Error("should have called delete once")
	}
}

func TestTimeRangeManagerDeleteNotFound(t *testing.T) {
	repoMock := &TimeRangeRepositoryMock{}

	manager := CreateTimeRangeManager(repoMock, DateMock{}).(*TimeRangeManagement)

	timeRange, createErr := manager.Create()

	if createErr != nil {
		t.Error("should not return an error")
	}

	errorFirstDelete := manager.Delete(timeRange.GetId())

	if errorFirstDelete != nil {
		t.Error("should not have returned an error")
	}

	errorSecondDelete := manager.Delete(timeRange.GetId())

	if !errors.Is(errorSecondDelete, TimeRangeManagerTimeRangeNotFoundErr) {
		t.Error("should have returned a not found error")
	}
}

func TestTimeRangeManagerDeleteError(t *testing.T) {
	deleteErr := errors.New("delete error")

	repoMock := &TimeRangeRepositoryMock{
		deleteError: deleteErr,
	}

	manager := CreateTimeRangeManager(repoMock, DateMock{}).(*TimeRangeManagement)

	timeRange, createErr := manager.Create()

	if createErr != nil {
		t.Error("should not return an error")
	}

	errorDelete := manager.Delete(timeRange.GetId())

	if !errors.Is(errorDelete, deleteErr) {
		t.Error("should have returned a delete error")
	}
}
