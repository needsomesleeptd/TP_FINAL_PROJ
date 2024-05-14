package models

import (
	"time"

	"github.com/google/uuid"
)

type Match struct {
	ID            uint64
	SessionID     uuid.UUID
	Datetime      time.Time
	GotFeedback   bool
	CardMatchedID uint64
}

