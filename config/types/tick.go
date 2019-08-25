package types

import "time"

type TickService struct {
	Interval time.Duration
	JobHandle func()
}
