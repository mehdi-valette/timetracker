package entity

import "time"

type Timestamp time.Time

func CreateTimestampSeconds(seconds int64) Timestamp {
	return Timestamp(time.Unix(seconds, 0))
}

func (t Timestamp) GetSeconds() int64 {
	return time.Time(t).Unix()
}

func (t Timestamp) IsBefore(t2 Timestamp) bool {
	return t.GetSeconds() < t2.GetSeconds()
}
