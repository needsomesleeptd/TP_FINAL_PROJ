package scroll_repo

import (
	"github.com/google/uuid"
	"test_backend_frontend/internal/models"
)

type ScrollRepository interface {
	AddScrollFact(fact *models.FactScrolled) error
	GetAllLikedPlaces(session_id uuid.UUID, user_id uint64) ([]uint64, error)
	GetAllUsersIdsForSession(session_id uuid.UUID) ([]uint64, error)
	GetAllPlaces(session_id uuid.UUID, user_id uint64) ([]uint64, error)
}
