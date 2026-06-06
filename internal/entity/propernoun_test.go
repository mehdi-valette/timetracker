package entity

import (
	"errors"
	"testing"

	"github.com/mehdi-valette/timetracker/internal/test"
)

func TestCreateName(t *testing.T) {
	name, err := CreateProperNoun("   hello   ")
	expected := "hello"

	if err != nil {
		t.Error(test.NoError(err))
	}

	if string(name) != expected {
		t.Errorf("Expected %s, got %s", expected, name)
	}
}

func TestCreateNameEmpty(t *testing.T) {
	name, err := CreateProperNoun("      ")

	if !errors.Is(err, EmptyNounErr) {
		t.Error("should have raised an error")
	}

	if string(name) != "" {
		t.Errorf("name should be empty, received: %s", name)
	}
}
