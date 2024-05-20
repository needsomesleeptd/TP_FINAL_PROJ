package models

import "time"

type Feedback struct {
	ID          uint
	UserID      uint64
	Description string
	HasGone     bool
	Datetime    time.Time
}
