package cache

import "time"

type stored struct {
	value   any
	validAt time.Time
}
