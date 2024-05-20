package models_dto // stands for data_transfer_objects

import (
	"time"

	"github.com/google/uuid"
)

type Match struct {
	ID            uint64    `json:"id"`
	SessionID     uuid.UUID `json:"session_id"`
	Datetime      time.Time `json:"datetime"`
	GotFeedback   bool      `json:"got_feedback"`
	CardMatchedID uint64    `json:"matched_card_id"`
	UserID        uint64    `json:"user_id"`
	MatchViewed   bool      `json:"match_viewed"`
}
