package entity

import "time"

type Dater interface {
	Now() Timestamp
}

type date struct{}

var _ Dater = date{}

func CreateDate() Dater {
	return date{}
}

func (d date) Now() Timestamp {
	return Timestamp(time.Now())
}
