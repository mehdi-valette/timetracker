package entity

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func createTestTask() (Task, *DaterMock) {
	mockDater := DaterMock{}

	return CreateTask(0, "my task", &mockDater), &mockDater
}

func TestCreateTask(t *testing.T) {
	task := CreateTask(1, "   test    ", Date{})

	if len(task.timeRanges) != 0 || task.name != "test" || task.id != 1 {
		t.Errorf("Shoudl have no time ranges, name=test and id=1, got %+v", task)
	}
}

func TestTaskIsEqual(t *testing.T) {
	type TestCase struct {
		task1    Task
		task2    Task
		expected bool
	}

	cases := []TestCase{
		{
			Task{
				id:   2,
				name: Name("test"),
			},
			Task{
				id:   3,
				name: Name("test2"),
			},
			false,
		},
		{
			Task{
				id:   2,
				name: Name("test"),
			},
			Task{
				id:   3,
				name: Name("test"),
			},
			false,
		},
		{
			Task{
				id:   2,
				name: Name("test"),
			},
			Task{
				id:   2,
				name: Name("test"),
			},
			true,
		},
	}

	for _, c := range cases {
		if (reflect.DeepEqual(c.task1, c.task2)) != c.expected {
			t.Errorf("equality of tasks %+v and %+v should be %t", c.task1, c.task2, c.expected)
		}
	}
}

func TestTaskRename(t *testing.T) {
	task := CreateTask(0, "     hello     ", Date{})

	if task.name != "hello" {
		t.Errorf("name should be 'hello")
	}

	task.Rename("world")

	if task.name != "world" {
		t.Errorf("name should be 'world'")
	}
}

func TestTaskSetTimeRange(t *testing.T) {
	task, _ := createTestTask()

	timeRange := CreateTimeRange(0, Date{})

	if len(task.timeRanges) != 0 {
		t.Errorf("shouldn't have time ranges yet")
	}

	task.SetTimeRange(&timeRange)

	if len(task.timeRanges) != 1 {
		t.Errorf("should have one time range")
	}
}

func TestTaskGetTimeRange(t *testing.T) {
	timeRange := CreateTimeRange(12, Date{})

	task, _ := createTestTask()

	task.SetTimeRange(&timeRange)

	firstGet, err := task.GetTimeRange(12)

	if err != nil {
		t.Errorf("should get a time range")
	}

	firstGet.End()

	secondGet, err := task.GetTimeRange(12)

	if err != nil {
		t.Errorf("should get a time range")
	}

	firstDuration, err1 := firstGet.Duration()
	secondDuration, err2 := secondGet.Duration()

	if err1 != nil || err2 != nil || firstDuration != secondDuration {
		t.Errorf("should be the same time range")
	}
}

func TestTaskGetTimeRangeNotFound(t *testing.T) {
	task, _ := createTestTask()

	_, err := task.GetTimeRange(0)

	if !errors.Is(err, TaskTimeRangeNotFoundErr) {
		t.Errorf("should return not found error")
	}
}

func TestTaskDurationNoTimeRange(t *testing.T) {
	task, _ := createTestTask()

	if duration, err := task.Duration(); err != nil || duration != 0 {
		t.Errorf("should return a duration of 0 when there's no time ranges")
	}
}

func TestTaskDurationSingleTimeRange(t *testing.T) {
	task, mockDate := createTestTask()
	currentTime := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	mockDate.Set(currentTime)

	timeRange := CreateTimeRange(0, mockDate)
	mockDate.Set(currentTime.Add(time.Hour))
	timeRange.End()

	task.SetTimeRange(&timeRange)

	if duration, err := task.Duration(); err != nil || duration != 3600 {
		t.Errorf("the task should last one hour.")
	}
}

func TestTaskDurationMultipleTimeRange(t *testing.T) {
	task, mockDate := createTestTask()

	currentTime := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	mockDate.Set(currentTime)

	// first time range
	firstTimeRange := CreateTimeRange(0, mockDate)
	mockDate.Set(currentTime.Add(time.Hour))
	firstTimeRange.End()

	// second time range
	secondTimeRange := CreateTimeRange(13, mockDate)
	mockDate.Set(currentTime.Add(3 * time.Hour))
	secondTimeRange.End()

	// join the time ranges to the task
	task.SetTimeRange(&firstTimeRange)
	task.SetTimeRange(&secondTimeRange)

	if duration, err := task.Duration(); err != nil || duration != 3600*3 {
		t.Errorf("the task should last three hour.")
	}
}

func TestTaskDurationUnfinishedTimeRange(t *testing.T) {
	task, mockDate := createTestTask()
	currentTime := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	mockDate.Set(currentTime)

	// first time range
	firstTimeRange := CreateTimeRange(12, mockDate)
	mockDate.Set(currentTime.Add(time.Hour))
	firstTimeRange.End()

	mockDate.Set(currentTime.Add(2 * time.Hour))

	// second time range
	secondTimeRange := CreateTimeRange(13, mockDate)

	mockDate.Set(currentTime.Add(3 * time.Hour))

	// join the time ranges to the task
	task.SetTimeRange(&firstTimeRange)
	task.SetTimeRange(&secondTimeRange)

	if duration, err := task.Duration(); err != nil || duration != 3600*2 {
		t.Errorf("the task should last two hour.")
	}
}

func TestTaskDurationTimeBackward(t *testing.T) {
	task, mockDate := createTestTask()
	currentTime := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	mockDate.Set(currentTime)

	timeRange := CreateTimeRange(12, mockDate)

	mockDate.Set(currentTime.Add(-time.Hour))

	task.SetTimeRange(&timeRange)

	if duration, err := task.Duration(); !errors.Is(err, TaskDurationImpossibleErr) || duration != 0 {
		t.Errorf("should return an error")
	}
}

func TestTaskDurationTwoUnfinishedTimeRange(t *testing.T) {
	task, mockDater := createTestTask()

	// first time range
	firstTimeRange := CreateTimeRange(12, mockDater)

	// second time range
	secondTimeRange := CreateTimeRange(13, mockDater)

	// join the time ranges to the task
	task.SetTimeRange(&firstTimeRange)
	task.SetTimeRange(&secondTimeRange)

	if duration, err := task.Duration(); !errors.Is(err, TaskTwoUnfinishedTimeRangeErr) || duration != 0 {
		t.Errorf("there should be an error")
	}
}

func TestIsRunningNoTimeRange(t *testing.T) {
	task, _ := createTestTask()

	if task.IsRunning() {
		t.Errorf("should not be running")
	}
}

func TestIsRunningSingleUnfinishedTimeRange(t *testing.T) {
	task, mockDater := createTestTask()

	timeRange := CreateTimeRange(12, mockDater)

	task.SetTimeRange(&timeRange)

	if !task.IsRunning() {
		t.Errorf("should be running")
	}
}

func TestIsRunningTwoFinishedTimeRanges(t *testing.T) {
	task, mockDater := createTestTask()

	firstTimeRange := CreateTimeRange(12, mockDater)
	mockDater.Set(mockDater.innerDate.Add(time.Hour))
	firstTimeRange.End()

	secondTimeRange := CreateTimeRange(13, mockDater)
	mockDater.Set(mockDater.innerDate.Add(time.Hour))
	secondTimeRange.End()

	task.SetTimeRange(&firstTimeRange)
	task.SetTimeRange(&secondTimeRange)

	if task.IsRunning() {
		t.Errorf("should not be running")
	}
}

func TestIsRunningOneFinishedOneUnfinishedimeRanges(t *testing.T) {
	task, mockDater := createTestTask()

	firstTimeRange := CreateTimeRange(12, mockDater)
	mockDater.Set(mockDater.innerDate.Add(time.Hour))
	firstTimeRange.End()

	secondTimeRange := CreateTimeRange(13, mockDater)

	task.SetTimeRange(&firstTimeRange)
	task.SetTimeRange(&secondTimeRange)

	if !task.IsRunning() {
		t.Errorf("should be running")
	}
}

func TestIsRunningTwoUnfinishedimeRanges(t *testing.T) {
	task, mockDater := createTestTask()

	firstTimeRange := CreateTimeRange(12, mockDater)

	secondTimeRange := CreateTimeRange(13, mockDater)

	task.SetTimeRange(&firstTimeRange)
	task.SetTimeRange(&secondTimeRange)

	if task.IsRunning() {
		t.Errorf("should not be running")
	}
}
