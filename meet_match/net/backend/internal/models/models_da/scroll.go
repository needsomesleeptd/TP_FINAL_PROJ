package models_da

import (
	"github.com/google/uuid"
	"test_backend_frontend/internal/models"
)

type FactScrolled struct {
	SessionID string `grom:"column:session_id"`
	UserID    uint64 `grom:"column:user_id"`
	PlacesID  uint64 `grom:"column:places_id"`
	IsLiked   bool   `grom:"column:is_liked"`
}

func (FactScrolled) TableName() string {
	return "fact_scrolled"
}

func ToPostgresFactScrolled(scrolled *models.FactScrolled) *FactScrolled {
	return &FactScrolled{
		SessionID: scrolled.SessionId.String(),
		UserID:    scrolled.UserId,
		PlacesID:  scrolled.PlacesId,
		IsLiked:   scrolled.IsLiked,
	}
}

func ToModelFactScrolled(scrolled *FactScrolled) *models.FactScrolled {
	// TODO: ой-ой-ой подумать
	session_id, _ := uuid.Parse(scrolled.SessionID)

	return &models.FactScrolled{
		SessionId: session_id,
		UserId:    scrolled.UserID,
		PlacesId:  scrolled.PlacesID,
		IsLiked:   scrolled.IsLiked,
	}
}
