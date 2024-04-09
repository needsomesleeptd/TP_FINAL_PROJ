package models

import "github.com/google/uuid"

type FactScrolled struct {
	SessionId uuid.UUID
	UserId    uint64
	PlacesId  uint64
	IsLiked   bool
}
