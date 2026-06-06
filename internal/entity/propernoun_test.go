package entity

import "testing"

func TestCreateName(t *testing.T) {
	name := CreateProperNoun("   hello   ")
	expected := "hello"

	if string(name) != expected {
		t.Errorf("Expected %s, got %s", expected, name)
	}
}
