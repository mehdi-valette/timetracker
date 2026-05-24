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

func CreateDateMock(innerDate time.Time) DateMock {
	return DateMock{
		innerDate: innerDate,
	}
}

type DateMock struct {
	innerDate time.Time
}

func (d DateMock) Now() Timestamp {
	return Timestamp(d.innerDate)
}

func (d *DateMock) Set(date time.Time) {
	d.innerDate = date
}

var _ Dater = DateMock{}
