package entity

import (
	"errors"
	"time"
)

var TimestampErrorNegativeDurationErr = errors.New("duration cannot be negative")

func CreateTimestampSeconds(seconds int64) Timestamp {
	return Timestamp(time.Unix(seconds, 0))
}

type Timestamper interface {
	GetSeconds() int64
	IsBefore(t2 Timestamp) bool
	TimeEllapsedSince(t2 Timestamp) (Duration, error)
}

type Timestamp time.Time

var _ Timestamper = Timestamp{}

func (t Timestamp) GetSeconds() int64 {
	return time.Time(t).Unix()
}

func (t Timestamp) IsBefore(t2 Timestamp) bool {
	return t.GetSeconds() < t2.GetSeconds()
}

func (t1 Timestamp) TimeEllapsedSince(t2 Timestamp) (Duration, error) {

	if t2.GetSeconds() > t1.GetSeconds() {
		return Duration(0), TimestampErrorNegativeDurationErr
	}

	return Duration(t1.GetSeconds() - t2.GetSeconds()), nil
}
