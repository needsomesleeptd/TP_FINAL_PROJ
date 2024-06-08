package models_dto // stands for data_transfer_objects

import (
	"time"
)

type Feedback struct {
	ID          uint      `json:"id"`
	UserID      uint64    `json:"user_id"`
	Description string    `json:"description,omitempty"`
	HasGone     bool      `json:"hasgone"`
	Datetime    time.Time `json:"datetime,omitempty"`
}
