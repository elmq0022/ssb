package timeutil

import (
	"time"
)

type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func (RealClock) Now() time.Time {
	return time.Now()
}

type FakeClock struct {
	FixedTime time.Time
}

func (fc FakeClock) Now() time.Time {
	return fc.FixedTime
}
