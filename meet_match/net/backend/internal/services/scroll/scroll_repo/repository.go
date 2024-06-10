package scroll_repo

import (
	"test_backend_frontend/internal/models"

	"github.com/google/uuid"
)

type ScrollRepository interface {
	AddScrollFact(fact *models.FactScrolled) error
	GetAllLikedPlaces(session_id uuid.UUID, user_id uint64) ([]uint64, error)
	GetAllUsersIdsForSession(session_id uuid.UUID) ([]uint64, error)
}
