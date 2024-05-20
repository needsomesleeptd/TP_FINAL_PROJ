package models_da

import (
	"test_backend_frontend/internal/models"
	"time"

	"github.com/google/uuid"
)

type Match struct {
	ID            uint64    `gorm:"primaryKey;column:id"`
	SessionID     uuid.UUID `gorm:"type:uuid;column:session_id"`
	Datetime      time.Time `gorm:"not null;column:datetime"`
	GotFeedback   bool      `gorm:"not null;column:got_feedback"`
	CardMatchedID uint64    `gorm:"not null;column:matched_card_id"`
	UserID        uint64    `gorm:"not null;column:user_id"`
	MatchViewed   bool      `gorm:"not null;column:match_viewed"`
}

func TODaMatch(match models.Match) Match {
	return Match(match)
}

func FromDaMatch(match Match) models.Match {
	return models.Match(match)
}

func FromDaMatchSlice(daMatches []Match) []models.Match {
	modelMatches := make([]models.Match, len(daMatches))
	for i, daMatch := range daMatches {
		modelMatches[i] = FromDaMatch(daMatch)
	}
	return modelMatches
}
