package entity

import (
	"testing"
	"time"
)

func TestDateNow(t *testing.T) {
	date := Date{}

	expected := Timestamp(time.Now())

	result := date.Now()

	if result.GetSeconds()-expected.GetSeconds() > 1 {
		t.Errorf("should have returned the current time")
	}
}
