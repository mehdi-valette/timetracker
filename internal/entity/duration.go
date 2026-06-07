package entity

import "fmt"

type Duration uint32

func (duration Duration) ToString() string {
	hours := duration / 3600
	minutes := (duration - (hours * 3600)) / 60
	seconds := duration - hours*3600 - minutes*60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}
