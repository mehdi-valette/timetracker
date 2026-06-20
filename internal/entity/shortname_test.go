package entity

import (
	"testing"
)

func TestCreateShortName(t *testing.T) {
	shortName := CreateShortName("  Ab28C dE ")
	expected := ShortName("ab28cde")

	if shortName != expected {
		t.Errorf("should have \"%s\", got \"%s\"", expected, shortName)
	}
}
