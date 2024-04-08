package models_da

import "github.com/google/uuid"

type FactScrolled struct {
	SessionID uuid.UUID `grom:"column:session_id"`
	UserID    uint64    `grom:"column:user_id"`
	PlacesID  uint64    `grom:"column:places_id"`
	IsLiked   bool      `grom:"column:is_liked"`
}

func (FactScrolled) TableName() string {
	return "fact_scrolled"
}
