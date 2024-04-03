package util

import "time"

func TimePtr(t time.Time) *time.Time {
	return &t
}

func PtrToTime(t *time.Time) time.Time {
	if t != nil {
		return *t
	}
	return time.Time{}
}
