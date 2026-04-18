package entity

import "time"

type Dater interface {
	Now() Timestamp
}

type Date struct{}

var _ Dater = Date{}

func (d Date) Now() Timestamp {
	return Timestamp(time.Now())
}
