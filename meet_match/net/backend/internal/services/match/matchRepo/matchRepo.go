package match_repo

import (
	"test_backend_frontend/internal/models"

	"github.com/google/uuid"
)

type IMatchRepo interface {
	GetMatchesBySession(sessionID uuid.UUID) ([]models.Match, error)
	GetMatchesNoFeedback(sessionID uuid.UUID) ([]models.Match, error)
	SaveMatch(match models.Match) error
}
