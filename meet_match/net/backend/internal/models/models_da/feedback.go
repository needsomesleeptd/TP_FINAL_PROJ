package models_da

import (
	"test_backend_frontend/internal/models"
	"time"
)

type Feedback struct {
	ID          uint      `gorm:"primaryKey;column:id"`
	UserID      uint64    `gorm:"not null;column:user_id"`
	Description string    `gorm:"column:description"`
	HasGone     bool      `gorm:"not null;column:has_gone"`
	Datetime    time.Time `gorm:"not null;column:datetime"`
}

func ToDAFeedback(feedack models.Feedback) Feedback {
	return Feedback(feedack)
}
