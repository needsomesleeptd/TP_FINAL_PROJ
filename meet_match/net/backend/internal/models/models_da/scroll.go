package models_da

import (
	"test_backend_frontend/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type FactScrolled struct {
	SessionID string    `gorm:"column:session_id"`
	UserID    uint64    `gorm:"column:user_id"`
	PlacesID  uint64    `gorm:"column:place_id"`
	IsLiked   bool      `gorm:"column:is_liked"`
	DateTime  time.Time `gorm:"column:datetime"`
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
		DateTime:  time.Now(),
	}
}

func ToModelFactScrolled(scrolled *FactScrolled) (*models.FactScrolled, error) {
	session_id, err := uuid.Parse(scrolled.SessionID)
	if err != nil {
		return nil, errors.Wrap(err, "ToModelFactScrolled parse error")
	}

	return &models.FactScrolled{
		SessionId: session_id,
		UserId:    scrolled.UserID,
		PlacesId:  scrolled.PlacesID,
		IsLiked:   scrolled.IsLiked,
	}, nil
}
