package entity

import (
	"testing"
)

func TestCreateTask(t *testing.T) {
	task := CreateTask(1, "   test    ")

	expected := Task{
		id:   1,
		name: Name("test"),
	}

	if task != expected {
		t.Errorf("Expected %v, got %v", expected, task)
	}
}

func TestIsEqual(t *testing.T) {
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
		if (c.task1 == c.task2) != c.expected {
			t.Errorf("equality of tasks %+v and %+v should be %t", c.task1, c.task2, c.expected)
		}
	}
}
