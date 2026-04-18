package entity

import "testing"

func TestCreateName(t *testing.T) {
	name := CreateName("   hello   ")
	expected := "hello"

	if string(name) != expected {
		t.Errorf("Expected %s, got %s", expected, name)
	}
}
